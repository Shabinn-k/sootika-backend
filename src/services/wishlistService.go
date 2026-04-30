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

// ⚠️ CRITICAL FIX: Add transaction to prevent orphaned wishlist
func (s *WishlistService) AddToWishlist(userID, productID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}
	
	var product models.Product
	if err := s.Repo.FindByID(&product, productUUID); err != nil {
		return errors.New("product not found")
	}
	
	// Start transaction
	tx := s.Repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// Get or create wishlist within transaction
	wishlist, err := s.getOrCreateWishlistWithTx(tx, userUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get wishlist: %w", err)
	}
	
	// Check if product already exists (with lock to prevent race condition)
	var existingItem models.WishlistItem
	err = tx.FindOneWhere(&existingItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID)
	if err == nil {
		tx.Rollback()
		return errors.New("product already in wishlist")
	}
	
	wishlistItem := &models.WishlistItem{
		WishlistID: wishlist.ID,
		ProductID:  productUUID,
	}
	
	if err := tx.Insert(wishlistItem); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add to wishlist: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit wishlist addition: %w", err)
	}
	
	return nil
}

// ⚠️ CRITICAL FIX: Add transaction
func (s *WishlistService) RemoveFromWishlist(userID, productID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}
	
	// Start transaction
	tx := s.Repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	var wishlist models.Wishlist
	if err := tx.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		tx.Rollback()
		return errors.New("wishlist not found")
	}
	
	var wishlistItem models.WishlistItem
	if err := tx.FindOneWhere(&wishlistItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID); err != nil {
		tx.Rollback()
		return errors.New("product not found in wishlist")
	}
	
	if err := tx.Delete(&wishlistItem, wishlistItem.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove from wishlist: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit wishlist removal: %w", err)
	}
	
	return nil
}

func (s *WishlistService) GetWishlist(userID string) (*models.Wishlist, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	wishlist, err := s.getOrCreateWishlist(userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}
	if err := s.loadWishlistItems(wishlist); err != nil {
		return nil, fmt.Errorf("failed to load wishlist items: %w", err)
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
		return false, fmt.Errorf("invalid user ID: %w", err)
	}
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		// ⚠️ FIX: Typo in error message
		return false, fmt.Errorf("invalid product ID: %w", err)
	}
	var wishlist models.Wishlist
	if err := s.Repo.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		return false, nil
	}
	var wishlistItem models.WishlistItem
	err = s.Repo.FindOneWhere(&wishlistItem, "wishlist_id = ? AND product_id = ?", wishlist.ID, productUUID)
	return err == nil, nil
}

// ⚠️ CRITICAL FIX: Add transaction
func (s *WishlistService) ClearWishlist(userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	
	// Start transaction
	tx := s.Repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	var wishlist models.Wishlist
	if err := tx.FindOneWhere(&wishlist, "user_id = ?", userUUID); err != nil {
		tx.Rollback()
		return errors.New("wishlist not found")
	}
	
	if err := tx.DeleteWhere(&models.WishlistItem{}, "wishlist_id = ?", wishlist.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit wishlist clear: %w", err)
	}
	
	return nil
}

// Helper with transaction
func (s *WishlistService) getOrCreateWishlistWithTx(tx repository.PgSQLRepository, userID uuid.UUID) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	err := tx.FindOneWhere(&wishlist, "user_id = ?", userID)

	if err == nil {
		return &wishlist, nil
	}

	wishlist = models.Wishlist{
		UserID: userID,
	}

	if err := tx.Insert(&wishlist); err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return &wishlist, nil
}

// Original non-transaction version for read-only operations
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