package controllers

import (
	"golang/src/models"
	"golang/src/services"
	"golang/utils/constant"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{
		Service: service,
	}
}

type CreateProductRequest struct {
	Title       string `json:"title" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	MainImage   string `json:"main_image" binding:"required"`
	SecondImage string `json:"second_image"`
	ThirdImage  string `json:"third_image"`
	Stock       int    `json:"stock" binding:"min=0"`
}

type UpdateProductRequest struct {
	Title       *string `json:"title"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	MainImage   *string `json:"main_image"`
	SecondImage *string `json:"second_image"`
	ThirdImage  *string `json:"third_image"`
	InStock     *bool   `json:"in_stock"`
	Stock       *int    `json:"stock"`
}

func (p *ProductController) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	product := &models.Product{
		Title:       req.Title,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		MainImage:   req.MainImage,
		SecondImage: req.SecondImage,
		ThirdImage:  req.ThirdImage,
		Stock:       req.Stock,
		InStock:     req.Stock > 0,
	}
	if err := p.Service.CreateProduct(product); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.CREATED, product)
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

func (p *ProductController) GetAllProduct(c *gin.Context) {
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
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	input := &services.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		MainImage:   req.MainImage,
		SecondImage: req.SecondImage,
		ThirdImage:  req.ThirdImage,
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

func (p *ProductController) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product Id is required"})
		return
	}
	if err := p.Service.DeleteProduct(productID); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Product deleted successfully"})
}

func (p *ProductController) UpdateProductStock(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Product ID is required"})
		return
	}
	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constant.BADREQUEST, gin.H{"error": err.Error()})
		return
	}
	if err := p.Service.UpdateProductStock(productID, req.Quantity); err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{"message": "Stock updated successfull"})
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
	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Search term is required"})
		return
	}
	products, err := p.Service.SearchProducts(searchTerm)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"data":  products,
		"count": len(products),
	})
}

func (p *ProductController) GetProductsByTitle(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(constant.BADREQUEST, gin.H{"error": "Title is required"})
		return
	}
	products, err := p.Service.GetProductsByTitle(title)
	if err != nil {
		c.JSON(constant.INTERNALSERVERERROR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constant.SUCCESS, gin.H{
		"data":  products,
		"count": len(products),
	})
}
