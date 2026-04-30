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
func (s *PaymentService) CreateOrder(amount int64, currency string, receipt string) (map[string]interface{}, error) {
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
        CreatedAt:       time.Now(),
    }
    
    if err := s.repo.Insert(paymentOrder); err != nil {
        return nil, fmt.Errorf("failed to save payment order: %w", err)
    }
    
    return order, nil
}

// VerifyPayment verifies the payment signature
func (s *PaymentService) VerifyPayment(orderID, paymentID, signature string) (bool, error) {
    var payment models.Payment
    if err := s.repo.FindOneWhere(&payment, "razorpay_order_id = ?", orderID); err != nil {
        return false, errors.New("payment order not found")
    }
    
    generatedSignature := s.generateSignature(orderID, paymentID)
    
    if generatedSignature != signature {
        return false, errors.New("signature verification failed")
    }
    
    updates := map[string]interface{}{
        "razorpay_payment_id": paymentID,
        "razorpay_signature":  signature,
        "status":              "paid",
        "updated_at":          time.Now(),
    }
    
    if err := s.repo.UpdateByFields(&models.Payment{}, payment.ID, updates); err != nil {
        return false, fmt.Errorf("failed to update payment status: %w", err)
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
    
    if expectedSignature != signature {
        return errors.New("invalid webhook signature")
    }
    
    var webhookData map[string]interface{}
    if err := json.Unmarshal(body, &webhookData); err != nil {
        return fmt.Errorf("failed to parse webhook: %w", err)
    }
    
    event, _ := webhookData["event"].(string)
    switch event {
    case "payment.captured":
        payload := webhookData["payload"].(map[string]interface{})
        payment := payload["payment"].(map[string]interface{})
        orderID := payment["order_id"].(string)
        paymentID := payment["id"].(string)
        
        var paymentRecord models.Payment
        if err := s.repo.FindOneWhere(&paymentRecord, "razorpay_order_id = ?", orderID); err == nil {
            updates := map[string]interface{}{
                "razorpay_payment_id": paymentID,
                "status":              "paid",
                "updated_at":          time.Now(),
            }
            s.repo.UpdateByFields(&models.Payment{}, paymentRecord.ID, updates)
        }
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

// RefundPayment initiates a refund - FIXED VERSION
func (s *PaymentService) RefundPayment(paymentID string, amount int64) (map[string]interface{}, error) {
    // Razorpay refund expects: (paymentID string, amount int, data map[string]interface{}, query map[string]string)
    refund, err := s.client.Payment.Refund(paymentID, int(amount), nil, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to refund payment: %w", err)
    }
    return refund, nil
}