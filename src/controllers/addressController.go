package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang/src/models"
    "golang/src/services"
    "golang/utils/constant"
)

type AddressController struct {
    Service *services.AddressService  // ⚠️ CHANGE THIS
}

func NewAddressController(service *services.AddressService) *AddressController {  // ⚠️ CHANGE THIS
    return &AddressController{Service: service}
}

func (c *AddressController) GetMyAddresses(ctx *gin.Context) {
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
    
    addresses, err := c.Service.GetUserAddresses(userUUID)
    if err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"addresses": addresses})
}

func (c *AddressController) AddAddress(ctx *gin.Context) {
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
    
    var address models.Address
    if err := ctx.ShouldBindJSON(&address); err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    if err := c.Service.AddAddress(userUUID, &address); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.CREATED, gin.H{"message": "Address added successfully", "address": address})
}

func (c *AddressController) UpdateAddress(ctx *gin.Context) {
    addressID := ctx.Param("id")
    addressUUID, err := uuid.Parse(addressID)
    if err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid address ID"})
        return
    }
    
    var updates map[string]interface{}
    if err := ctx.ShouldBindJSON(&updates); err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
        return
    }
    
    if err := c.Service.UpdateAddress(addressUUID, updates); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Address updated successfully"})
}

func (c *AddressController) DeleteAddress(ctx *gin.Context) {
    addressID := ctx.Param("id")
    addressUUID, err := uuid.Parse(addressID)
    if err != nil {
        ctx.JSON(constant.BADREQUEST, gin.H{"error": "Invalid address ID"})
        return
    }
    
    if err := c.Service.DeleteAddress(addressUUID); err != nil {
        ctx.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(constant.SUCCESS, gin.H{"message": "Address deleted successfully"})
}