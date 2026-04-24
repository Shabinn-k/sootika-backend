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

func (a *AdminController) Dashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var totalProducts int64
	var totalUsers int64

	a.repo.Count(&models.Product{}, &totalProducts)
	a.repo.Count(&models.User{}, &totalUsers)

	c.JSON(constant.SUCCESS, gin.H{
		"message":  "Welcome to Admin Dashboard",
		"admin_id": userID,
		"role":     role,
		"stats": gin.H{
			"total_products": totalProducts,
			"total_users":    totalUsers,
		},
	})
}

func (a *AdminController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := a.repo.FindAll(&users); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	c.JSON(constant.SUCCESS, gin.H{
		"data":  users,
		"count": len(users),
	})
}

func (a *AdminController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	var user models.User
	if err := a.repo.FindByID(&user, userID); err != nil {
		c.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	user.Password = ""
	c.JSON(constant.SUCCESS, user)
}

func (a *AdminController) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	var req struct {
		Role string `json:"role" binding:"required,oneof=user admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{
		"role": req.Role,
	}
	if err := a.repo.UpdateByFields(&models.User{}, userID, updates); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to update user role"})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message": "User role updated successfully",
		"role":    req.Role,
	})
}

func (a *AdminController) ToggleBlockUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	currentAdminID, _ := c.Get("user_id")
	if userID == currentAdminID {
		c.JSON(constant.BADREQUEST, gin.H{"error": "You cannot block/unblock yourself"})
		return
	}
	var user models.User
	if err := a.repo.FindByID(&user, userID); err != nil {
		c.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	newBlockStatus := !user.IsBlocked
	updates := map[string]interface{}{
		"is_blocked": newBlockStatus,
	}
	if err := a.repo.UpdateByFields(&models.User{}, userID, updates); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to toggle user block status"})
		return
	}
	message := "User unblocked successfully"
	if newBlockStatus {
		message = "User blocked successfully"
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message":    message,
		"is_blocked": newBlockStatus,
	})
}

func (a *AdminController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "User ID is required"})
		return
	}
	currentAdminID, _ := c.Get("user_id")
	if userID == currentAdminID {
		c.JSON(constant.BADREQUEST, gin.H{"error": "You cannot delete yourself"})
		return
	}
	var user models.User
	if err := a.repo.FindByID(&user, userID); err != nil {
		c.JSON(constant.NOTFOUND, gin.H{"error": "User not found"})
		return
	}
	if err := a.repo.Delete(&user, userID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "User deleted successfully"})
}

func (a *AdminController) GetTotalProducts(c *gin.Context) {
	var count int64
	a.repo.Count(&models.Product{}, &count)
	c.JSON(constant.SUCCESS, gin.H{"total_products": count})
}

func (a *AdminController) GetTotalUsers(c *gin.Context) {
	var count int64
	a.repo.Count(&models.User{}, &count)
	c.JSON(constant.SUCCESS, gin.H{"total_users": count})
}
