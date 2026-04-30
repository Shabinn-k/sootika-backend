package services

import (
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "golang/src/models"
    "golang/src/repository"
)

type OrderService struct {
    repo repository.PgSQLRepository
}

func NewOrderService(repo repository.PgSQLRepository) *OrderService {
    return &OrderService{repo: repo}
}
type OrderItemInput struct {
    ProductID uuid.UUID
    Quantity  int
}
// CreateOrder creates a new order
func (s *OrderService) CreateOrder(userID uuid.UUID, items []OrderItemInput, addressID uuid.UUID, paymentMethod string) (*models.Order, error) {
    var user models.User
    if err := s.repo.FindByID(&user, userID); err != nil {
        return nil, fmt.Errorf("user not found")
    }
    
    var address models.Address
    if err := s.repo.FindByID(&address, addressID); err != nil {
        return nil, fmt.Errorf("address not found")
    }
    
    var subtotal int64
    var orderItems []models.OrderItem
    
    for _, item := range items {
        var product models.Product
        if err := s.repo.FindByID(&product, item.ProductID); err != nil {
            return nil, fmt.Errorf("product not found: %v", item.ProductID)
        }
        
        total := product.Price * int64(item.Quantity)
        subtotal += total
        
        orderItems = append(orderItems, models.OrderItem{
            ProductID: item.ProductID,
            Title:     product.Title,
            Name:      product.Name,
            Image:     product.MainImage,
            Price:     product.Price,
            Quantity:  item.Quantity,
            Total:     total,
        })
    }
    
    shippingCost := int64(80)
    tax := int64(0)
    discount := int64(0)
    total := subtotal + shippingCost + tax - discount
    
    orderNumber := fmt.Sprintf("SOO%d%d", time.Now().UnixNano(), userID[:8])
    orderNumber = orderNumber[:20]
    
    order := &models.Order{
        OrderNumber:    orderNumber,
        UserID:         userID,
        Total:          total,
        Subtotal:       subtotal,
        ShippingCost:   shippingCost,
        Tax:            tax,
        Discount:       discount,
        PaymentMethod:  paymentMethod,
        PaymentStatus:  "pending",
        OrderStatus:    "pending",
        Track:          "Pending",
        ShippingAddress: address,
    }
    
    if err := s.repo.Insert(order); err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    for i := range orderItems {
        orderItems[i].OrderID = order.ID
        if err := s.repo.Insert(&orderItems[i]); err != nil {
            return nil, fmt.Errorf("failed to create order items: %w", err)
        }
    }
    
    order.Items = orderItems
    return order, nil
}

// GetUserOrders gets all orders for a user
func (s *OrderService) GetUserOrders(userID uuid.UUID) ([]models.Order, error) {
    var orders []models.Order
    if err := s.repo.FindAllWhere(&orders, "user_id = ?", userID); err != nil {
        return nil, err
    }
    
    for i := range orders {
        var items []models.OrderItem
        s.repo.FindAllWhere(&items, "order_id = ?", orders[i].ID)
        orders[i].Items = items
    }
    
    return orders, nil
}

// GetOrderByID gets order by ID
func (s *OrderService) GetOrderByID(orderID string) (*models.Order, error) {
    var order models.Order
    if err := s.repo.FindByID(&order, orderID); err != nil {
        return nil, err
    }
    
    var items []models.OrderItem
    s.repo.FindAllWhere(&items, "order_id = ?", order.ID)
    order.Items = items
    
    return &order, nil
}

// UpdateOrderStatus updates order status
func (s *OrderService) UpdateOrderStatus(orderID string, status string) error {
    updates := map[string]interface{}{
        "order_status": status,
        "track":        status,
        "updated_at":   time.Now(),
    }
    return s.repo.UpdateByFields(&models.Order{}, orderID, updates)
}

// UpdatePaymentStatus updates payment status
func (s *OrderService) UpdatePaymentStatus(orderID string, paymentStatus string, paymentID string) error {
    updates := map[string]interface{}{
        "payment_status":      paymentStatus,
        "razorpay_payment_id": paymentID,
        "updated_at":          time.Now(),
    }
    return s.repo.UpdateByFields(&models.Order{}, orderID, updates)
}