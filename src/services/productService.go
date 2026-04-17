package services

import (
	"errors"
	"fmt"
	"golang/src/models"
	"golang/src/repository"
)

type ProductService struct {
	Repo repository.PgSQLRepository
}

func SetupProductService(repo repository.PgSQLRepository)*ProductService{
	return &ProductService{Repo: repo}
}

type UpdateProductInput struct{
	Name *string
	Description *string
	Price *int64
	MainImage *string
	SecondImage *string
	ThirdImage *string
	InStock *bool
	Stock *int
}

func (s *ProductService)CreateProduct(product *models.Product)error{
	if product==nil{
		return errors.New("Product data in nil") 
	}
	if product.Title==""{
		return errors.New("Product title required")
	}
	if product.Name==""{
		return errors.New("Product Name is required")
	}
	if product.Price<=0{
		return errors.New("Product Price must be Greater than 0")
	}
	if err:=s.Repo.Insert(product);err!=nil{
		return fmt.Errorf("failed to create product: %w",err)
	}
	return nil
}

func (s *ProductService)GetAllProducts()([]models.Product,error){
	var products []models.Product
	if err:=s.Repo.FindAll(&products);err!=nil{
		return nil,fmt.Errorf("Failed to fetch products: %w",err)
	}
	return products,nil
}

func (s *ProductService)GetProductByID(productID string)(*models.Product,error){
	var product models.Product
	if err:=s.Repo.FindByID(&product,"id = ?",productID);err!=nil{

		return nil,fmt.Errorf("Product not fount: %w",err)
	}
	return &product,nil
}

func (s *ProductService)UpdateProduct(productID string,input *UpdateProductInput)(*models.Product,error){
	var product models.Product
	if err:=s.Repo.FindByID(&product,"id = ?",productID);err!=nil{
		return nil,fmt.Errorf("Product not found: %w",err)
	}
	updates:=map[string]interface{}{}

	if input.Name!=nil{
		updates["name"]=*input.Name
	}
	if input.Description!=nil{
		updates["description"]=*input.Description
	}
	if input.Price!=nil{
		if *input.Price<=0{
			return nil,errors.New("Price must be greater than zero")
		}
		updates["price"]=*input.Price
	}
	if input.MainImage!=nil{
		updates["main_image"]=*input.MainImage
	}
	if input.SecondImage!=nil{
		updates["second_image"]=*input.SecondImage
	}
	if input.ThirdImage!=nil{
		updates["third_image"]=*input.ThirdImage
	}
	if input.InStock!=nil{
		updates["in_stock"]=*input.InStock
	}
	if input.Stock!=nil{
		if *input.Stock<0{
			return nil,errors.New("Stock cannot be negative")
		}
		updates["stock"]=*input.Stock	
		updates["in_stock"]=*input.Stock>0
	}
	if len(updates)>0{
		if err:=s.Repo.UpdateByFields(&models.Product{},productID,updates);err!=nil{
			return nil,fmt.Errorf("Failed to update product: %w",err)
		}
	}
	if err:=s.Repo.FindByID(&product,"id = ?",productID);err!=nil{
		return nil,fmt.Errorf("Failed to fetch updated product: %w",err)
	}
	return &product,nil
}

func (s *ProductService)DeleteProduct(productID string)error{
	var product models.Product
	if err:=s.Repo.FindByID(&product,"id = ?",productID);err!=nil{
		return fmt.Errorf("Product not found: %w",err)
	}
	if err:=s.Repo.Delete(&product,productID);err!=nil{
		return fmt.Errorf("Failed to delete product: %w",err)
	}
	return nil
}

func (s *ProductService)UpdateProductStock(productID string,quantity int)error{
	product,err:=s.GetProductByID(productID)
	if err!=nil{
		return err
	}
	newStock:=product.Stock-quantity
	if newStock<0{
		return errors.New("insufficient stock")
	}
	updates:=map[string]interface{}{
		"stock":newStock,
		"in_stock":newStock>0,
	}
	if err:=s.Repo.UpdateByFields(&models.Product{},productID,updates);err!=nil{
		return fmt.Errorf("Failed to update stock: %w",err)
	}
	return nil
}

func (s *ProductService)GetInStock()([]models.Product,error){
	var product []models.Product
	if err:=s.Repo.FindWhere(&product,"in_stock = ?",true);err!=nil{
		return nil,fmt.Errorf("Failed to fetch in-stock products: %w",err)
	}
	return product,nil
}

func (s *ProductService)SearchProducts(searchTerm string)([]models.Product,error){
	var products []models.Product
	query:="name ILIKE ? OR description ILIKE ?"
	searchPattern:="%"+searchTerm+"%"
	if err:=s.Repo.FindWhere(&products,query,searchPattern,searchPattern);err!=nil{
		return nil,fmt.Errorf("Failed to search products: %w",err)
	}
	return products,nil
}

func (s *ProductService)GetProductsByTitle(title string)([]models.Product,error){
	var products []models.Product
	if err:=s.Repo.FindWhere(&products,"title ILIKE ?","%"+title+"%");err!=nil{
		return nil,fmt.Errorf("Failed to fetch products by title: %w",err)
	}
	return products,nil
}