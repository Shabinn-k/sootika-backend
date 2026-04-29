package controllers

import (
	"github.com/gin-gonic/gin"
	"golang/src/services"
	"golang/utils/constant"
	"golang/utils/logger"
	"golang/utils/validation"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{authService: service}
}

type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=6,password"`
}

type VerifyOTP struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (a *AuthController) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	logger.Log.Info("Signup:", req.Email)
	err := a.authService.Signup(req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		logger.Log.Error("Signup failed:", err)
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message": "Signup successful, OTP sent",
	})
}

func (a *AuthController) VerifyOTP(c *gin.Context) {
	var req VerifyOTP
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	err := a.authService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		logger.Log.Warn("OTP failed:", err)
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message": "Account verified",
	})
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (c *AuthController) ResendOTP(ctx *gin.Context) {
	var req ResendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	if err := c.authService.ResendOTP(req.Email); err != nil {
		ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(constant.SUCCESS, gin.H{"message": "OTP resent successfully"})
}

func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	access, refresh, user, err := a.authService.Login(req.Email, req.Password)
	if err != nil {
		logger.Log.Error("Login failed:", err)
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"role":          user.Role,
	})
}

func (a *AuthController) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	newAccess, newRefresh, err := a.authService.Refresh(body.RefreshToken)
	if err != nil {
		logger.Log.Warn("Refresh failed", err)
		c.JSON(constant.UNAUTHORIZED, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"new_access_token":  newAccess,
		"new_refresh_token": newRefresh,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(constant.BADREQUEST, validation.FormatValidationErrors(err))
		return
	}
	err := a.authService.Logout(body.RefreshToken)
	if err != nil {
		logger.Log.Error("Logout failed:", err)
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Logged out successfully"})
}

func (a *AuthController) Dashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	user, err := a.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "failed to get user info"})
		return
	}
	if role == "admin" {
		var totalProducts int64
		var totalUsers int64
		a.authService.GetDashboardStats(&totalProducts, &totalUsers)
		c.JSON(constant.SUCCESS, gin.H{
			"message":  "Welcome to Admin Dashboard",
			"admin_id": userID,
			"role":     role,
			"stats": gin.H{
				"total_products": totalProducts,
				"total_users":    totalUsers,
			},
		})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"message": "Welcome to User Dashboard",
		"user_id": userID,
		"name":    user.Name,
		"role":    role,
	})
}
