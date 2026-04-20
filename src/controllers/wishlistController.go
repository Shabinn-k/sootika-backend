package controllers

import (
	"golang/src/services"
	"golang/utils/constant"

	"github.com/gin-gonic/gin"
)

type WishlistController struct {
	Service *services.WishlistService
}

func NewWishlistController(service *services.WishlistService) *WishlistController {
	return &WishlistController{
		Service: service,
	}
}

func (w *WishlistController) GetWishlist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	wishlist, err := w.Service.GetWishlist(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message":  "Wishlist fetched successfully",
		"wishlist": wishlist,
		"items":    len(wishlist.Items),
	})
}
func (w *WishlistController) AddToWishlist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not Authentication"})
		return
	}
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	if err := w.Service.AddToWishlist(userID.(string), req.ProductID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.CREATED, gin.H{"message": "Product added to wishlist successfully"})
}

func (w *WishlistController) RemoveFromWishlist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}
	if err := w.Service.RemoveFromWishlist(userID.(string), productID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Product removed from wishlist successfully"})
}

func (w *WishlistController) IsInWishlist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}
	isInWishlist, err := w.Service.IsInWishlist(userID.(string), productID)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"is_in_wishlist": isInWishlist})
}
func (w *WishlistController) GetWishlistCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not Authenticated"})
		return
	}
	count, err := w.Service.GetWishlistCount(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"count": count})
}

func (w *WishlistController) ClearWishlist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
		return
	}
	if err := w.Service.ClearWishlist(userID.(string)); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Wishlist cleared successfully"})
}
