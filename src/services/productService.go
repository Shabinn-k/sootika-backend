package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang/src/models"
	"golang/src/repository"
	"golang/utils/uploads"
	"mime/multipart"
)

type ProductService struct {
	Repo repository.PgSQLRepository
}

func NewProductService(repo repository.PgSQLRepository) *ProductService {
	return &ProductService{Repo: repo}
}

type CreateProductInput struct {
	Title       string
	Name        string
	Description string
	Price       int64
	Stock       int
	MainImage   *multipart.FileHeader
	SecondImage *multipart.FileHeader
	ThirdImage  *multipart.FileHeader
}
type UpdateProductInput struct {
	Title       *string
	Name        *string
	Description *string
	Price       *int64
	MainImage   *multipart.FileHeader
	SecondImage *multipart.FileHeader
	ThirdImage  *multipart.FileHeader
	InStock     *bool
	Stock       *int
}

func (s *ProductService) CreateProduct(input *CreateProductInput) (*models.Product, error) {
	if input == nil {
		return nil, errors.New("Product data in nil")
	}
	if input.Title == "" {
		return nil, errors.New("Product title required")
	}
	if input.Name == "" {
		return nil, errors.New("Product Name is required")
	}
	if input.Price <= 0 {
		return nil, errors.New("Product Price must be Greater than 0")
	}
	if input.MainImage == nil {
		return nil, errors.New("Main image is required")
	}
	mainFile, err := input.MainImage.Open()
	if err != nil {
		return nil, fmt.Errorf("Failed to open main image: %w", err)
	}
	defer mainFile.Close()
	mainResult, err := uploads.UploadImageFile(mainFile, input.MainImage.Filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to upload main image: %w", err)
	}
	var secondResult *uploads.CloudinaryResult
	if input.SecondImage != nil {
		secondFile, err := input.SecondImage.Open()
		if err == nil {
			defer secondFile.Close()
			secondResult, _ = uploads.UploadImageFile(secondFile, input.SecondImage.Filename)
		}
	}
	var thirdResult *uploads.CloudinaryResult
	if input.ThirdImage != nil {
		thirdFile, err := input.ThirdImage.Open()
		if err == nil {
			defer thirdFile.Close()
			thirdResult, _ = uploads.UploadImageFile(thirdFile, input.ThirdImage.Filename)
		}
	}

	product := &models.Product{
		Title:             input.Title,
		Name:              input.Name,
		Description:       input.Description,
		Price:             input.Price,
		Stock:             input.Stock,
		InStock:           input.Stock > 0,
		MainImage:         mainResult.URL,
		MainImagePublicID: mainResult.PublicID,
	}
	if secondResult != nil {
		product.SecondImage = secondResult.URL
		product.SecondImagePublicID = secondResult.PublicID
	}
	if thirdResult != nil {
		product.ThirdImage = thirdResult.URL
		product.ThirdImagePublicID = thirdResult.PublicID
	}
	if err := s.Repo.Insert(product); err != nil {
		uploads.DeleteImage(mainResult.PublicID)
		if secondResult != nil {
			uploads.DeleteImage(secondResult.PublicID)
		}
		if thirdResult != nil {
			uploads.DeleteImage(thirdResult.PublicID)
		}
		return nil, fmt.Errorf("Failed to create products: %w", err)
	}
	return product, nil
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.Repo.FindAll(&products); err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	return products, nil
}

func (s *ProductService) GetProductByID(productID string) (*models.Product, error) {
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product id: %w", err)
	}
	var product models.Product
	if err := s.Repo.FindByID(&product, productUUID); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return &product, nil
}

func (s *ProductService) UpdateProduct(productID string, input *UpdateProductInput) (*models.Product, error) {
	var product models.Product
	if err := s.Repo.FindByID(&product, productID); err != nil {
		return nil, fmt.Errorf("Product not found: %w", err)
	}
	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Price != nil {
		if *input.Price <= 0 {
			return nil, errors.New("Price must be greater than zero")
		}
		updates["price"] = *input.Price
	}
	if input.InStock != nil {
		updates["in_stock"] = *input.InStock
	}
	if input.Stock != nil {
		if *input.Stock < 0 {
			return nil, errors.New("Stock cannot be negative")
		}
		updates["stock"] = *input.Stock
		updates["in_stock"] = *input.Stock > 0
	}
	if len(updates) > 0 {
		if err := s.Repo.UpdateByFields(&models.Product{}, productID, updates); err != nil {
			return nil, fmt.Errorf("Failed to update product: %w", err)
		}
	}
	if err := s.Repo.FindByID(&product, productID); err != nil {
		return nil, fmt.Errorf("Failed to fetch updated product: %w", err)
	}
	return &product, nil
}
func (s *ProductService) UpdateProductImage(productID string, imageType string, newImage *multipart.FileHeader) error {
	var product models.Product
	if err := s.Repo.FindByID(&product, productID); err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	file, err := newImage.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	result, err := uploads.UploadImageFile(file, newImage.Filename)
	if err != nil {
		return err
	}
	var oldPublicID string
	switch imageType {
	case "main":
		oldPublicID = product.MainImagePublicID
		product.MainImage = result.URL
		product.MainImagePublicID = result.PublicID
	case "second":
		oldPublicID = product.SecondImagePublicID
		product.SecondImage = result.URL
		product.SecondImagePublicID = result.PublicID
	case "third":
		oldPublicID = product.ThirdImagePublicID
		product.ThirdImage = result.URL
		product.ThirdImagePublicID = result.PublicID
	default:
		return errors.New("invalid image type")
	}
	if oldPublicID != "" {
		uploads.DeleteImage(oldPublicID)
	}
	updates := map[string]interface{}{
		imageType + "_image":           result.URL,
		imageType + "_image_public_id": result.PublicID,
	}
	return s.Repo.UpdateByFields(&models.Product{}, productID, updates)
}

func (s *ProductService) DeleteProduct(productID string) error {
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product id: %w", err)
	}
	var product models.Product
	if err := s.Repo.FindByID(&product, productUUID); err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	if product.MainImagePublicID != "" {
		uploads.DeleteImage(product.MainImagePublicID)
	}
	if product.SecondImagePublicID != "" {
		uploads.DeleteImage(product.SecondImagePublicID)
	}
	if product.ThirdImagePublicID != "" {
		uploads.DeleteImage(product.ThirdImagePublicID)
	}
	if err := s.Repo.Delete(&product, productUUID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (s *ProductService) UpdateProductStock(productID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	product, err := s.GetProductByID(productID)
	if err != nil {
		return err
	}
	newStock := product.Stock - quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}
	productUUID, _ := uuid.Parse(productID)
	updates := map[string]interface{}{
		"stock":    newStock,
		"in_stock": newStock > 0,
	}
	if err := s.Repo.UpdateByFields(&models.Product{}, productUUID, updates); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

func (s *ProductService) GetInStockProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.Repo.FindAllWhere(&products, "in_stock = ?", true); err != nil {
		return nil, fmt.Errorf("failed to fetch in-stock products: %w", err)
	}
	return products, nil
}

func (s *ProductService) SearchProducts(term string) ([]models.Product, error) {
	var products []models.Product
	searchPattern := "%" + term + "%"
	query := `name ILIKE ? OR description ILIKE ? OR title ILIKE ?`
	if err := s.Repo.FindAllWhere(&products, query, searchPattern, searchPattern, searchPattern); err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	return products, nil
}
