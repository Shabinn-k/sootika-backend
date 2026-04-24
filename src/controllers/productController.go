package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang/src/services"
	"golang/utils/constant"
	"strconv"
)

type ProductController struct {
	Service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{
		Service: service,
	}
}

func (p *ProductController) CreateProduct(c *gin.Context) {
	title := c.PostForm("title")
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	stockStr := c.PostForm("stock")

	if title == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "title is required"})
		return
	}
	if name == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "name is required"})
		return
	}
	if description == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "description is required"})
		return
	}
	if priceStr == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "price is required"})
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": "invalid price format"})
		return
	}
	stock := 0
	if stockStr != "" {
		stock, err = strconv.Atoi(stockStr)
		if err != nil {
			c.JSON(constant.BADREQUEST, gin.H{"error": "invalid stock format"})
			return
		}
	}

	mainImage, err := c.FormFile("main_image")
	if err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": "main_image is required"})
		return
	}

	secondImage, _ := c.FormFile("second_image")
	thirdImage, _ := c.FormFile("third_image")

	input := &services.CreateProductInput{
		Title:       title,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		MainImage:   mainImage,
		SecondImage: secondImage,
		ThirdImage:  thirdImage,
	}

	product, err := p.Service.CreateProduct(input)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.CREATED, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (p *ProductController) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}
	product, err := p.Service.GetProductByID(productID)
	if err != nil {
		c.JSON(constant.NOTFOUND, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, product)
}

func (p *ProductController) GetAllProducts(c *gin.Context) {
	products, err := p.Service.GetAllProducts()
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"data":  products,
		"count": len(products),
	})
}

func (p *ProductController) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}
	var req struct {
		Title       *string `json:"title"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Price       *int64  `json:"price"`
		InStock     *bool   `json:"in_stock"`
		Stock       *int    `json:"stock"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}

	input := &services.UpdateProductInput{
		Title:       req.Title,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		InStock:     req.InStock,
		Stock:       req.Stock,
	}

	product, err := p.Service.UpdateProduct(productID, input)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.SUCCESS, product)
}

func (p *ProductController) UpdateProductImage(c *gin.Context) {
	productID := c.Param("id")
	imageType := c.Param("type")

	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}

	if imageType != "main" && imageType != "second" && imageType != "third" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "invalid image type. Use: main, second, or third"})
		return
	}

	newImage, err := c.FormFile("image")
	if err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": "image file is required"})
		return
	}

	if err := p.Service.UpdateProductImage(productID, imageType, newImage); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.SUCCESS, gin.H{
		"message": fmt.Sprintf("%s image updated successfully", imageType),
	})
}

func (p *ProductController) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}

	if err := p.Service.DeleteProduct(productID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.SUCCESS, gin.H{"message": "Product deleted successfully"})
}

func (p *ProductController) GetInStockProducts(c *gin.Context) {
	products, err := p.Service.GetInStockProducts()
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.SUCCESS, gin.H{
		"data":  products,
		"count": len(products),
	})
}

func (p *ProductController) SearchProducts(c *gin.Context) {
	term := c.Query("q")
	if term == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "search query required"})
		return
	}
	products, err := p.Service.SearchProducts(term)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}

	c.JSON(constant.SUCCESS, gin.H{
		"data":  products,
		"count": len(products),
	})
}
