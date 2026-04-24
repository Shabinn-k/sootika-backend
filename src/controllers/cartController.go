package controllers

import (
	"github.com/gin-gonic/gin"
	"golang/src/services"
	"golang/utils/constant"
)

type CartController struct {
	Services *services.CartService
}

func NewCartController(service *services.CartService) *CartController {
	return &CartController{
		Services: service,
	}
}

func (s *CartController) GetCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	cart, err := s.Services.GetCart(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message": "Cart fetched successfully",
		"cart":    cart,
		"Items":   len(cart.Items),
	})
}

func (s *CartController) AddToCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User is not authenticated"})
		return
	}
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	if err := s.Services.AddToCart(userID.(string), req.ProductID, req.Quantity); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.CREATED, gin.H{"message": "Product added to cart successfully"})
}

func (s *CartController) UpdateCartItemQuantity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	cartItemID := c.Param("item_id")
	if cartItemID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Cart item ID is required"})
		return
	}
	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	if err := s.Services.UpdateCartItemQuantity(userID.(string), cartItemID, req.Quantity); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Cart item quantity updated successfully"})
}
func (s *CartController) RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User is not authenticated"})
		return
	}
	cartItemID := c.Param("item_id")
	if cartItemID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Cart item ID is required"})
		return
	}
	if err := s.Services.RemoveFromCart(userID.(string), cartItemID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Product removed from cart successfully"})
}

func (s *CartController) GetCartCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User is not authenticated"})
		return
	}
	count, err := s.Services.GetCartCount(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"count": count})
}

func (s *CartController) GetCartTotal(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User is not authenticated"})
		return
	}
	total, err := s.Services.GetCartTotal(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"total": total})
}

func (s *CartController) ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	if err := s.Services.ClearCart(userID.(string)); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Cart cleared successfully"})
}
