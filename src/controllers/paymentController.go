package controllers

import (
    "io"
    "github.com/gin-gonic/gin"
    "golang/src/services"
    "golang/utils/constant"
)

type PaymentController struct {
    paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
    return &PaymentController{
        paymentService: paymentService,
    }
}

// Renamed to RazorpayOrderRequest to avoid conflict
type RazorpayOrderRequest struct {
    Amount   int64  `json:"amount" binding:"required"`
    Currency string `json:"currency" binding:"required"`
    Receipt  string `json:"receipt" binding:"required"`
}

// CreateOrder - POST /api/payment/create-order
func (p *PaymentController) CreateOrder(c *gin.Context) {
    var req RazorpayOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    order, err := p.paymentService.CreateOrder(req.Amount, req.Currency, req.Receipt)
    if err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(constant.SUCCESS, gin.H{
        "order_id":  order["id"],
        "amount":    order["amount"],
        "currency":  order["currency"],
    })
}

type VerifyPaymentRequest struct {
    OrderID   string `json:"order_id" binding:"required"`
    PaymentID string `json:"payment_id" binding:"required"`
    Signature string `json:"signature" binding:"required"`
}

// VerifyPayment - POST /api/payment/verify
func (p *PaymentController) VerifyPayment(c *gin.Context) {
    var req VerifyPaymentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    verified, err := p.paymentService.VerifyPayment(req.OrderID, req.PaymentID, req.Signature)
    if err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    if !verified {
        c.JSON(constant.BADREQUEST, gin.H{"error": "Payment verification failed"})
        return
    }
    
    c.JSON(constant.SUCCESS, gin.H{
        "message": "Payment verified successfully",
    })
}

// Webhook - POST /api/payment/webhook
func (p *PaymentController) Webhook(c *gin.Context) {
    signature := c.GetHeader("X-Razorpay-Signature")
    
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": "Failed to read body"})
        return
    }
    
    if err := p.paymentService.HandleWebhook(body, signature); err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(constant.SUCCESS, gin.H{"status": "ok"})
}