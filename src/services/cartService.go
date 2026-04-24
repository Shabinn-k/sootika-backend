package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang/src/models"
	"golang/src/repository"
)

type CartService struct {
	Repo repository.PgSQLRepository
}

func NewCartService(repo repository.PgSQLRepository) *CartService {
	return &CartService{
		Repo: repo,
	}
}

func (s *CartService) getOrCreateCart(userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := s.Repo.FindOneWhere(&cart, "user_id = ?", userID)
	if err == nil {
		return &cart, nil
	}
	cart = models.Cart{
		UserID: userID,
	}
	if err := s.Repo.Insert(&cart); err != nil {
		return nil, fmt.Errorf("Failed to create cart: %w", err)
	}
	return &cart, nil
}

func (s *CartService) AddToCart(userID, productID string, quantity int) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product id: %w", err)
	}
	if quantity <= 0 {
		quantity = 1
	}
	var product models.Product
	if err := s.Repo.FindByID(&product, productUUID); err != nil {
		return errors.New("product not found")
	}
	cart, err := s.getOrCreateCart(userUUID)
	if err != nil {
		return err
	}
	var item models.CartItem
	err = s.Repo.FindOneWhere(&item,
		"cart_id = ? AND product_id = ?",
		cart.ID, productUUID,
	)
	newQuantity := quantity

	if err == nil {
		newQuantity = item.Quantity + quantity
	}
	if newQuantity > product.Stock {
		return fmt.Errorf("only %d items available in stock", product.Stock)
	}
	if err == nil {
		return s.Repo.UpdateByFields(&models.CartItem{}, item.ID, map[string]interface{}{
			"quantity": item.Quantity + quantity,
		})
	}
	newItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: productUUID,
		Quantity:  quantity,
		Price:     int(product.Price),
	}
	return s.Repo.Insert(&newItem)
}

func (s *CartService) UpdateCartItemQuantity(userID, cartItemID string, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("Invalid user ID: %w", err)
	}
	itemUUID, err := uuid.Parse(cartItemID)
	if err != nil {
		return fmt.Errorf("Invalid cart item ID: %w", err)
	}
	var cart models.Cart
	if err := s.Repo.FindOneWhere(&cart, "user_id = ?", userUUID); err != nil {
		return errors.New("Cart not found")
	}
	var cartItem models.CartItem
	if err := s.Repo.FindOneWhere(&cartItem, "id = ? AND cart_id = ?", itemUUID, cart.ID); err != nil {
		return errors.New("Cart item not found")
	}
	var product models.Product
	if err := s.Repo.FindByID(&product, cartItem.ProductID); err != nil {
		return errors.New("Product not found")
	}
	if quantity > product.Stock {
		return fmt.Errorf("only %d items available in stock", product.Stock)
	}
	return s.Repo.UpdateByFields(&models.CartItem{}, cartItem.ID, map[string]interface{}{
		"quantity": quantity,
	})
}

func (s *CartService) RemoveFromCart(userID, cartItemID string) error {
	userUUID, _ := uuid.Parse(userID)
	itemUUID, _ := uuid.Parse(cartItemID)

	var cart models.Cart
	if err := s.Repo.FindOneWhere(&cart, "user_id = ?", userUUID); err != nil {
		return errors.New("cart not found")
	}

	return s.Repo.GetDB().
		Where("id = ? AND cart_id = ?", itemUUID, cart.ID).
		Delete(&models.CartItem{}).Error
}

func (s *CartService) loadCartItems(cart *models.Cart) error {
	var items []models.CartItem

	err := s.Repo.GetDB().
		Where("cart_id = ?", cart.ID).
		Preload("Product").
		Find(&items).Error

	if err != nil {
		return err
	}

	cart.Items = items
	return nil
}
func (s *CartService) GetCart(userID string) (*models.Cart, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	cart, err := s.getOrCreateCart(userUUID)
	if err != nil {
		return nil, err
	}

	if err := s.loadCartItems(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) GetCartCount(userID string) (int, error) {
	userUUID, _ := uuid.Parse(userID)

	var cart models.Cart
	if err := s.Repo.FindOneWhere(&cart, "user_id = ?", userUUID); err != nil {
		return 0, nil
	}

	var count int64
	s.Repo.GetDB().
		Model(&models.CartItem{}).
		Where("cart_id = ?", cart.ID).
		Count(&count)

	return int(count), nil
}

func (s *CartService) GetCartTotal(userID string) (int64, error) {
	cart, err := s.GetCart(userID)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, item := range cart.Items {
		total += int64(item.Price * item.Quantity)
	}

	return total, nil
}

func (s *CartService) ClearCart(userID string) error {
	userUUID, _ := uuid.Parse(userID)

	var cart models.Cart
	if err := s.Repo.FindOneWhere(&cart, "user_id = ?", userUUID); err != nil {
		return errors.New("cart not found")
	}

	return s.Repo.GetDB().
		Where("cart_id = ?", cart.ID).
		Delete(&models.CartItem{}).Error
}
