package services

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "time"
    
    "github.com/razorpay/razorpay-go"
    "github.com/google/uuid"
    "golang/src/models"
    "golang/src/repository"
    "golang/config"
)

type PaymentService struct {
    repo      repository.PgSQLRepository
    client    *razorpay.Client
    keySecret string
}

func NewPaymentService(repo repository.PgSQLRepository, cfg *config.Config) *PaymentService {
    client := razorpay.NewClient(cfg.Razorpay.KeyID, cfg.Razorpay.KeySecret)
    
    return &PaymentService{
        repo:      repo,
        client:    client,
        keySecret: cfg.Razorpay.KeySecret,
    }
}

// CreateOrder creates a Razorpay order
func (s *PaymentService) CreateOrder(amount int64, currency string, receipt string, userID uuid.UUID) (map[string]interface{}, error) {
    if amount <= 0 {
        return nil, errors.New("invalid amount")
    }
    
    data := map[string]interface{}{
        "amount":          amount,
        "currency":        currency,
        "receipt":         receipt,
        "payment_capture": 1,
    }
    
    order, err := s.client.Order.Create(data, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    paymentOrder := &models.Payment{
        ID:              uuid.New(),
        OrderID:         receipt,
        RazorpayOrderID: order["id"].(string),
        Amount:          amount,
        Currency:        currency,
        Status:          "created",
        UserID:          userID,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }
    
    if err := s.repo.Insert(paymentOrder); err != nil {
        return nil, fmt.Errorf("failed to save payment order: %w", err)
    }
    
    return order, nil
}

// VerifyPayment verifies payment signature and updates order
func (s *PaymentService) VerifyPayment(orderID, paymentID, signature string, orderUUID uuid.UUID, userID uuid.UUID) (bool, error) {
    // Start transaction
    tx := s.repo.BeginTransaction()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    var payment models.Payment
    if err := tx.FindOneWhere(&payment, "razorpay_order_id = ? AND user_id = ?", orderID, userID); err != nil {
        tx.Rollback()
        return false, errors.New("payment order not found")
    }
    
    generatedSignature := s.generateSignature(orderID, paymentID)
    
    if generatedSignature != signature {
        tx.Rollback()
        return false, errors.New("signature verification failed")
    }
    
    // Update payment record
    updates := map[string]interface{}{
        "razorpay_payment_id": paymentID,
        "razorpay_signature":  signature,
        "status":              "paid",
        "updated_at":          time.Now(),
    }
    
    if err := tx.UpdateByFields(&models.Payment{}, payment.ID, updates); err != nil {
        tx.Rollback()
        return false, fmt.Errorf("failed to update payment status: %w", err)
    }
    
    // Update actual order payment status
    orderUpdates := map[string]interface{}{
        "payment_status":      "paid",
        "razorpay_payment_id": paymentID,
        "updated_at":          time.Now(),
    }
    
    if err := tx.UpdateByFields(&models.Order{}, orderUUID, orderUpdates); err != nil {
        tx.Rollback()
        return false, fmt.Errorf("failed to update order payment status: %w", err)
    }
    
    if err := tx.Commit(); err != nil {
        return false, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return true, nil
}

func (s *PaymentService) generateSignature(orderID, paymentID string) string {
    data := orderID + "|" + paymentID
    h := hmac.New(sha256.New, []byte(s.keySecret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

// HandleWebhook processes Razorpay webhook
func (s *PaymentService) HandleWebhook(body []byte, signature string) error {
    secret := s.keySecret
    expectedSignature := s.generateWebhookSignature(body, secret)
    
    if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
        return errors.New("invalid webhook signature")
    }
    
    var webhookData map[string]interface{}
    if err := json.Unmarshal(body, &webhookData); err != nil {
        return fmt.Errorf("failed to parse webhook: %w", err)
    }
    
    event, ok := webhookData["event"].(string)
    if !ok {
        return errors.New("invalid webhook event")
    }
    
    switch event {
    case "payment.captured":
        payload, ok := webhookData["payload"].(map[string]interface{})
        if !ok {
            return errors.New("invalid webhook payload")
        }
        
        payment, ok := payload["payment"].(map[string]interface{})
        if !ok {
            return errors.New("invalid payment data")
        }
        
        orderID, ok := payment["order_id"].(string)
        if !ok {
            return errors.New("invalid order ID")
        }
        
        paymentID, ok := payment["id"].(string)
        if !ok {
            return errors.New("invalid payment ID")
        }
        
        // Use transaction for webhook
        tx := s.repo.BeginTransaction()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
            }
        }()
        
        var paymentRecord models.Payment
        if err := tx.FindOneWhere(&paymentRecord, "razorpay_order_id = ?", orderID); err == nil {
            updates := map[string]interface{}{
                "razorpay_payment_id": paymentID,
                "status":              "paid",
                "updated_at":          time.Now(),
            }
            tx.UpdateByFields(&models.Payment{}, paymentRecord.ID, updates)
            
            // Update order if linked
            if paymentRecord.OrderID != "" {
                orderUpdates := map[string]interface{}{
                    "payment_status": "paid",
                    "updated_at":     time.Now(),
                }
                tx.UpdateByFields(&models.Order{}, paymentRecord.OrderID, orderUpdates)
            }
        }
        
        tx.Commit()
    }
    
    return nil
}

func (s *PaymentService) generateWebhookSignature(body []byte, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write(body)
    return hex.EncodeToString(h.Sum(nil))
}

// FetchPayment fetches payment details from Razorpay
func (s *PaymentService) FetchPayment(paymentID string) (map[string]interface{}, error) {
    payment, err := s.client.Payment.Fetch(paymentID, nil, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch payment: %w", err)
    }
    return payment, nil
}

// RefundPayment initiates a refund
func (s *PaymentService) RefundPayment(paymentID string, amount int64) (map[string]interface{}, error) {
    if amount <= 0 {
        return nil, errors.New("invalid refund amount")
    }
    
    refund, err := s.client.Payment.Refund(paymentID, int(amount), nil, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to refund payment: %w", err)
    }
    return refund, nil
}