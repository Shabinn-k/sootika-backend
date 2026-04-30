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
    
    userUUID, _ := uuid.Parse(userID.(string))
    addressUUID, _ := uuid.Parse(req.AddressID)
    
    var items []services.OrderItemInput
    for _, item := range req.Items {
        productUUID, _ := uuid.Parse(item.ProductID)
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
    
    userUUID, _ := uuid.Parse(userID.(string))
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

// GetOrderByID - GET /api/orders/:id
func (c *OrderController) GetOrderByID(ctx *gin.Context) {
    orderID := ctx.Param("id")
    order, err := c.orderService.GetOrderByID(orderID)
    if err != nil {
        ctx.JSON(constant.NOTFOUND, gin.H{"error": "Order not found"})
        return
    }
    
    ctx.JSON(constant.SUCCESS, order)
}

// CancelOrder - PUT /api/orders/:id/cancel
func (c *OrderController) CancelOrder(ctx *gin.Context) {
    orderID := ctx.Param("id")
    
    if err := c.orderService.UpdateOrderStatus(orderID, "cancelled"); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Order cancelled successfully"})
}