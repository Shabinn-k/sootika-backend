package controllers

import (
    "github.com/gin-gonic/gin"
    "golang/src/models"
    "golang/src/repository"
    "golang/src/services"
    "golang/utils/constant"
    "time"
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
    role, exists := c.Get("role")
    if !exists || role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    userID, _ := c.Get("user_id")

    var totalProducts int64
    var totalUsers int64
    var pendingFeedback int64
    var recentUsers []models.User
    var totalOrders int64
    var pendingOrders int64
    var totalRevenue int64

    if err := a.repo.Count(&models.Product{}, &totalProducts); err != nil {
        totalProducts = 0
    }
    if err := a.repo.Count(&models.User{}, &totalUsers); err != nil {
        totalUsers = 0
    }
    
    if err := a.repo.GetDB().Model(&models.Feedback{}).Where("feed = ?", "pending").Count(&pendingFeedback).Error; err != nil {
        pendingFeedback = 0
    }
    
    a.repo.GetDB().Model(&models.Order{}).Count(&totalOrders)
    a.repo.GetDB().Model(&models.Order{}).Where("order_status = ?", "pending").Count(&pendingOrders)
    
    if err := a.repo.GetDB().Model(&models.Order{}).
        Where("order_status = ?", "delivered").
        Select("COALESCE(SUM(total), 0)").
        Scan(&totalRevenue).Error; err != nil {
        totalRevenue = 0
    }

    a.repo.GetDB().Order("created_at desc").Limit(5).Find(&recentUsers)
    
    for i := range recentUsers {
        recentUsers[i].Password = ""
    }

    c.JSON(constant.SUCCESS, gin.H{
        "message":  "Welcome to Admin Dashboard",
        "admin_id": userID,
        "role":     role,
        "stats": gin.H{
            "total_products":   totalProducts,
            "total_users":      totalUsers,
            "pending_feedback": pendingFeedback,
            "total_revenue":    totalRevenue,
            "recent_users":     recentUsers,
            "total_orders":     totalOrders,
            "pending_orders":   pendingOrders,
        },
    })
}

func (a *AdminController) GetAllUsers(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    var users []models.User
    if err := a.repo.GetDB().Order("created_at desc").Find(&users).Error; err != nil {
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
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
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
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
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
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
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
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
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
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    var count int64
    a.repo.Count(&models.Product{}, &count)
    c.JSON(constant.SUCCESS, gin.H{"total_products": count})
}

func (a *AdminController) GetTotalUsers(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    var count int64
    a.repo.Count(&models.User{}, &count)
    c.JSON(constant.SUCCESS, gin.H{"total_users": count})
}

func (a *AdminController) GetAllFeedbacks(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    var feedbacks []models.Feedback
    if err := a.repo.FindAll(&feedbacks); err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to fetch feedbacks"})
        return
    }
    c.JSON(constant.SUCCESS, feedbacks)
}

func (a *AdminController) ApproveFeedback(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    feedbackID := c.Param("id")
    if feedbackID == "" {
        c.JSON(constant.BADREQUEST, gin.H{"error": "Feedback ID is required"})
        return
    }
    updates := map[string]interface{}{
        "feed": "approved",
    }
    if err := a.repo.UpdateByFields(&models.Feedback{}, feedbackID, updates); err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to approve feedback"})
        return
    }
    c.JSON(constant.SUCCESS, gin.H{"message": "Feedback approved successfully"})
}

func (a *AdminController) DeleteFeedback(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    feedbackID := c.Param("id")
    if feedbackID == "" {
        c.JSON(constant.BADREQUEST, gin.H{"error": "Feedback ID is required"})
        return
    }
    var feedback models.Feedback
    if err := a.repo.FindByID(&feedback, feedbackID); err != nil {
        c.JSON(constant.NOTFOUND, gin.H{"error": "Feedback not found"})
        return
    }
    if err := a.repo.Delete(&feedback, feedbackID); err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to delete feedback"})
        return
    }
    c.JSON(constant.SUCCESS, gin.H{"message": "Feedback deleted successfully"})
}

func (a *AdminController) GetAllOrders(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    var orders []models.Order
    
    if err := a.repo.GetDB().
        Preload("Items").
        Order("created_at desc").
        Find(&orders).Error; err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    for i := range orders {
        var user models.User
        if err := a.repo.FindByID(&user, orders[i].UserID); err == nil {
            user.Password = ""
            orders[i].User = user
        }
    }
    
    c.JSON(constant.SUCCESS, gin.H{
        "data":  orders,
        "count": len(orders),
    })
}

func (a *AdminController) UpdateOrderStatus(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "admin" {
        c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Admin access required"})
        return
    }
    
    orderID := c.Param("id")
    if orderID == "" {
        c.JSON(constant.BADREQUEST, gin.H{"error": "Order ID is required"})
        return
    }
    
    var req struct {
        Status string `json:"status" binding:"required,oneof=Pending Confirmed Shipped Delivered Cancelled"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    updates := map[string]interface{}{
        "track":        req.Status,
        "order_status": req.Status,
        "updated_at":   time.Now(),
    }
    
    if err := a.repo.UpdateByFields(&models.Order{}, orderID, updates); err != nil {
        c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": "Failed to update order status"})
        return
    }
    
    c.JSON(constant.SUCCESS, gin.H{
        "message": "Order status updated successfully",
        "status":  req.Status,
    })
}