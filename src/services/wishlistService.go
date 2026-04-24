package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang/src/models"
	"golang/src/repository"
)

type WishlistService struct {
	Repo repository.PgSQLRepository
}

func NewWishlistService(repo repository.PgSQLRepository) *WishlistService {
	return &WishlistService{
		Repo: repo,
	}
}

func (s *WishlistService) AddToWishlist(userID, productID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("Invalid user ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("Invalid product ID: %w", err)
	}
	var product models.Product
	if err := s.Repo.FindByID(&product, productUUID); err != nil {
		return errors.New("Product not found")
	}
	wishlist, err := s.getOrCreateWishlist(userUUID)
	if err != nil {
		return fmt.Errorf("Failed to get wishlist: %w", err)
	}
	var existingItem models.WishlistItem
	err = s.Repo.FindOneWhere(&existingItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID)
	if err == nil {
		return errors.New("Product already in wishlist")
	}
	wishlistItem := &models.WishlistItem{
		WishlistID: wishlist.ID,
		ProductID:  productUUID,
	}
	if err := s.Repo.Insert(wishlistItem); err != nil {
		return fmt.Errorf("Failed to add to wishlist: %w", err)
	}
	return nil
}

func (s *WishlistService) RemoveFromWishlist(userID, productID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("Invalid product ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("Invalid product ID: %w", err)
	}
	var wishlist models.Wishlist
	if err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		return errors.New("Wishlist not found")
	}
	var wishlistItem models.WishlistItem
	if err := s.Repo.FindOneWhere(&wishlistItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID); err != nil {
		return errors.New("Product not found in wishlist")
	}
	if err := s.Repo.Delete(&wishlistItem, wishlistItem.ID); err != nil {
		return fmt.Errorf("Failed to remove from wishlist: %w", err)
	}
	return nil
}

func (s *WishlistService) GetWishlist(userID string) (*models.Wishlist, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("Invalid user ID: %w", err)
	}
	wishlist, err := s.getOrCreateWishlist(userUUID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get wishlist: %w", err)
	}
	if err := s.loadWishlistItems(wishlist); err != nil {
		return nil, fmt.Errorf("Failed to load wishlist items: %w", err)
	}
	return wishlist, nil
}

func (s *WishlistService) GetWishlistCount(userID string) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}

	var wishlist models.Wishlist
	if err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		return 0, nil
	}

	var count int64
	if err := s.Repo.GetDB().
		Model(&models.WishlistItem{}).
		Where("wishlist_id = ?", wishlist.ID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count wishlist items: %w", err)
	}

	return int(count), nil
}

func (s *WishlistService) IsInWishlist(userID, productID string) (bool, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("Invalid user ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return false, fmt.Errorf("Invalid product IF: %w", err)
	}
	var wishlist models.Wishlist
	if err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		return false, nil
	}
	var wishlistItem models.WishlistItem
	err = s.Repo.FindOneWhere(&wishlistItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID)
	return err == nil, nil
}

func (s *WishlistService) ClearWishlist(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("Invalid user ID: %w", err)
	}
	var wishlist models.Wishlist
	if err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		return errors.New("Wishlist not found")
	}
	if err := s.Repo.DeleteWhere(&models.WishlistItem{}, "wishlist_id = ?", wishlist.ID); err != nil {
		return fmt.Errorf("Failed to clear wishlist: %w", err)
	}
	return nil
}

func (s *WishlistService) getOrCreateWishlist(userID uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userID)

	if err == nil {
		return &wishlist, nil
	}

	wishlist = models.Wishlist{
		UserID: userID,
	}

	if err := s.Repo.Insert(&wishlist); err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return &wishlist, nil
}

func (s *WishlistService) loadWishlistItems(wishlist *models.Wishlist) error {
	var items []models.WishlistItem
	if err := s.Repo.GetDB().
		Where("wishlist_id = ?", wishlist.ID).
		Preload("Product").
		Find(&items).Error; err != nil {
		return err
	}
	wishlist.Items = items
	return nil
}
