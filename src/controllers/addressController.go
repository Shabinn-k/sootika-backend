package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang/src/models"
    "golang/src/repository"
    "golang/utils/constant"
)

type AddressController struct {
    repo repository.PgSQLRepository
}

func NewAddressController(repo repository.PgSQLRepository) *AddressController {
    return &AddressController{repo: repo}
}

type AddressInput struct {
    Name      string `json:"name"`
    Address   string `json:"address" binding:"required"`
    City      string `json:"city" binding:"required"`
    State     string `json:"state" binding:"required"`
    Pincode   string `json:"pincode" binding:"required"`
    Phone     string `json:"phone"`
    IsDefault bool   `json:"is_default"`
}

// GetMyAddresses - GET /api/addresses
func (c *AddressController) GetMyAddresses(ctx *gin.Context) {
    userID, _ := ctx.Get("user_id")
    userUUID, _ := uuid.Parse(userID.(string))
    
    var addresses []models.Address
    if err := c.repo.FindAllWhere(&addresses, "user_id = ?", userUUID); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{
        "data":  addresses,
        "count": len(addresses),
    })
}

// AddAddress - POST /api/addresses
func (c *AddressController) AddAddress(ctx *gin.Context) {
    userID, _ := ctx.Get("user_id")
    userUUID, _ := uuid.Parse(userID.(string))
    
    var req AddressInput
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    // If this is default, remove default from other addresses
    if req.IsDefault {
        c.repo.UpdateByFields(&models.Address{}, nil, map[string]interface{}{
            "is_default": false,
        })
    }
    
    address := &models.Address{
        UserID:    userUUID,
        Name:      req.Name,
        Address:   req.Address,
        City:      req.City,
        State:     req.State,
        Pincode:   req.Pincode,
        Phone:     req.Phone,
        IsDefault: req.IsDefault,
    }
    
    if err := c.repo.Insert(address); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.CREATED, gin.H{
        "message": "Address added successfully",
        "address": address,
    })
}

// UpdateAddress - PUT /api/addresses/:id
func (c *AddressController) UpdateAddress(ctx *gin.Context) {
    addressID := ctx.Param("id")
    userID, _ := ctx.Get("user_id")
    
    var req AddressInput
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    // Verify address belongs to user
    var address models.Address
    if err := c.repo.FindOneWhere(&address, "id = ? AND user_id = ?", addressID, userID); err != nil {
        ctx.JSON(constant.NOTFOUND, gin.H{"error": "Address not found"})
        return
    }
    
    updates := map[string]interface{}{
        "address":    req.Address,
        "city":       req.City,
        "state":      req.State,
        "pincode":    req.Pincode,
        "is_default": req.IsDefault,
        "updated_at": "now()",
    }
    
    if req.Name != "" {
        updates["name"] = req.Name
    }
    if req.Phone != "" {
        updates["phone"] = req.Phone
    }
    
    if req.IsDefault {
        c.repo.UpdateByFields(&models.Address{}, nil, map[string]interface{}{
            "is_default": false,
        })
    }
    
    if err := c.repo.UpdateByFields(&models.Address{}, addressID, updates); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Address updated successfully"})
}

// DeleteAddress - DELETE /api/addresses/:id
func (c *AddressController) DeleteAddress(ctx *gin.Context) {
    addressID := ctx.Param("id")
    userID, _ := ctx.Get("user_id")
    
    var address models.Address
    if err := c.repo.FindOneWhere(&address, "id = ? AND user_id = ?", addressID, userID); err != nil {
        ctx.JSON(constant.NOTFOUND, gin.H{"error": "Address not found"})
        return
    }
    
    if err := c.repo.Delete(&address, addressID); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Address deleted successfully"})
}