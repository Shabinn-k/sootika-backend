package controllers

import (
	"github.com/gin-gonic/gin"
	"golang/src/models"
	"golang/src/repository"
	"golang/src/services"
	"golang/utils/constant"
)

type AdminController struct {
	productService *services.ProductService
	repo           repository.PgSQLRepository
}

func NewAdminController(
	productService *services.ProductService,
	repo repository.PgSQLRepository,
) *AdminController {
	return &AdminController{
		productService: productService,
		repo:           repo,
	}
}

func (c *AdminController) Dashboard(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	role, _ := ctx.Get("role")

	var totalProducts int64
	var totalUsers int64

	c.repo.Count(&models.Product{}, &totalProducts)
	c.repo.Count(&models.User{}, &totalUsers)

	ctx.JSON(constant.SUCCESS, gin.H{
		"message":  "Welcome to Admin Dashboard",
		"admin_id": userID,
		"role":     role,
		"stats": gin.H{
			"total_products": totalProducts,
			"total_users":    totalUsers,
		},
	})
}

func (c *AdminController) GetAllUsers(ctx *gin.Context) {
	var users []models.User
	if err := c.repo.FindAll(&users); err != nil {
		ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	ctx.JSON(constant.SUCCESS, gin.H{
		"data":  users,
		"count": len(users),
	})
}

func (c *AdminController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	var user models.User
	if err := c.repo.FindByID(&user, userID); err != nil {
		ctx.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	user.Password = ""
	ctx.JSON(constant.SUCCESS, user)
}

func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{
		"role": req.Role,
	}
	if err := c.repo.UpdateByFields(&models.User{}, userID, updates); err != nil {
		ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to update user role"})
		return
	}
	ctx.JSON(constant.SUCCESS, gin.H{
		"message": "User role updated successfully",
		"role":    req.Role,
	})
}

func (c *AdminController) ToggleBlockUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	currentAdminID, _ := ctx.Get("user_id")
	if userID == currentAdminID {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "You cannot block/unblock yourself"})
		return
	}
	var user models.User
	if err := c.repo.FindByID(&user, userID); err != nil {
		ctx.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	newBlockStatus := !user.IsBlocked
	updates := map[string]interface{}{
		"is_blocked": newBlockStatus,
	}
	if err := c.repo.UpdateByFields(&models.User{}, userID, updates); err != nil {
		ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to toggle user block status"})
		return
	}
	message := "User unblocked successfully"
	if newBlockStatus {
		message = "User blocked successfully"
	}
	ctx.JSON(constant.SUCCESS, gin.H{
		"message":    message,
		"is_blocked": newBlockStatus,
	})
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	currentAdminID, _ := ctx.Get("user_id")
	if userID == currentAdminID {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": "You cannot delete yourself"})
		return
	}
	var user models.User
	if err := c.repo.FindByID(&user, userID); err != nil {
		ctx.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	if err := c.repo.Delete(&user, userID); err != nil {
		ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(constant.SUCCESS, gin.H{"message": "User deleted successfully"})
}

func (c *AdminController) GetTotalProducts(ctx *gin.Context) {
	var count int64
	c.repo.Count(&models.Product{}, &count)
	ctx.JSON(constant.SUCCESS, gin.H{"total_products": count})
}

func (c *AdminController) GetTotalUsers(ctx *gin.Context) {
	var count int64
	c.repo.Count(&models.User{}, &count)
	ctx.JSON(constant.SUCCESS, gin.H{"total_users": count})
}
