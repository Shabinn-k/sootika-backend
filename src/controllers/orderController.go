package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang/src/services"
    "golang/utils/constant"
)

type OrderController struct {
    orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
    return &OrderController{orderService: orderService}
}

type OrderItemInput struct {
    ProductID string `json:"product_id" binding:"required"`
    Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
    Items         []OrderItemInput `json:"items" binding:"required"`
    AddressID     string           `json:"address_id" binding:"required"`
    PaymentMethod string           `json:"payment_method" binding:"required"`
}

// CreateOrder - POST /api/orders
func (c *OrderController) CreateOrder(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
        return
    }
    
    var req CreateOrderRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    // ⚠️ FIX: Validate UUID parsing
    userUUID, err := uuid.Parse(userID.(string))
    if err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid user ID"})
        return
    }
    
    addressUUID, err := uuid.Parse(req.AddressID)
    if err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid address ID"})
        return
    }
    
    var items []services.OrderItemInput
    for _, item := range req.Items {
        productUUID, err := uuid.Parse(item.ProductID)
        if err != nil {
            ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid product ID: " + item.ProductID})
            return
        }
        items = append(items, services.OrderItemInput{
            ProductID: productUUID,
            Quantity:  item.Quantity,
        })
    }
    
    order, err := c.orderService.CreateOrder(userUUID, items, addressUUID, req.PaymentMethod)
    if err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.CREATED, gin.H{
        "message": "Order created successfully",
        "order":   order,
    })
}

// GetMyOrders - GET /api/orders
func (c *OrderController) GetMyOrders(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
        return
    }
    
    userUUID, err := uuid.Parse(userID.(string))
    if err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid user ID"})
        return
    }
    
    orders, err := c.orderService.GetUserOrders(userUUID)
    if err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{
        "data":  orders,
        "count": len(orders),
    })
}

// ⚠️ CRITICAL FIX: Add ownership check
func (c *OrderController) GetOrderByID(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
        return
    }
    
    orderID := ctx.Param("id")
    order, err := c.orderService.GetOrderByID(orderID)
    if err != nil {
        ctx.JSON(constant.NOTFOUND, gin.H{"error": "Order not found"})
        return
    }
    
    // ⚠️ FIX: Verify order belongs to user
    userUUID, _ := uuid.Parse(userID.(string))
    if order.UserID != userUUID {
        ctx.JSON(constant.FORBIDDEN, gin.H{"error": "Access denied: order does not belong to you"})
        return
    }
    
    ctx.JSON(constant.SUCCESS, order)
}

// ⚠️ CRITICAL FIX: Add ownership check
func (c *OrderController) CancelOrder(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
        return
    }
    
    orderID := ctx.Param("id")
    
    // ⚠️ FIX: Get order first to check ownership
    order, err := c.orderService.GetOrderByID(orderID)
    if err != nil {
        ctx.JSON(constant.NOTFOUND, gin.H{"error": "Order not found"})
        return
    }
    
    // Verify order belongs to user
    userUUID, _ := uuid.Parse(userID.(string))
    if order.UserID != userUUID {
        ctx.JSON(constant.FORBIDDEN, gin.H{"error": "Access denied: cannot cancel another user's order"})
        return
    }
    
    // Check if order can be cancelled
    if order.OrderStatus == "cancelled" {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Order already cancelled"})
        return
    }
    
    if order.OrderStatus == "delivered" {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Cannot cancel delivered order"})
        return
    }
    
    if err := c.orderService.UpdateOrderStatus(orderID, "cancelled"); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Order cancelled successfully"})
}