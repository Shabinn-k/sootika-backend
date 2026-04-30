package services

import (
    "github.com/google/uuid"
    "golang/src/models"
    "golang/src/repository"
)

type AddressService struct {
    repo repository.PgSQLRepository
}

func NewAddressService(repo repository.PgSQLRepository) *AddressService {
    return &AddressService{repo: repo}
}

func (s *AddressService) GetUserAddresses(userID uuid.UUID) ([]models.Address, error) {
    var addresses []models.Address
    if err := s.repo.FindAllWhere(&addresses, "user_id = ?", userID); err != nil {
        return nil, err
    }
    return addresses, nil
}

func (s *AddressService) AddAddress(userID uuid.UUID, address *models.Address) error {
    address.UserID = userID
    if address.IsDefault {
        // Remove default from other addresses of this user
        if err := s.repo.UpdateByFieldsWhere(&models.Address{}, 
            map[string]interface{}{"is_default": false}, 
            "user_id = ?", userID); err != nil {
            return err
        }
    }
    return s.repo.Insert(address)
}

func (s *AddressService) UpdateAddress(addressID uuid.UUID, updates map[string]interface{}) error {
    // First get the address to know which user it belongs to
    var address models.Address
    if err := s.repo.FindByID(&address, addressID); err != nil {
        return err
    }
    
    // If setting as default, remove default from user's other addresses
    if isDefault, ok := updates["is_default"]; ok && isDefault == true {
        if err := s.repo.UpdateByFieldsWhere(&models.Address{}, 
            map[string]interface{}{"is_default": false}, 
            "user_id = ? AND id != ?", address.UserID, addressID); err != nil {
            return err
        }
    }
    
    return s.repo.UpdateByFields(&models.Address{}, addressID, updates)
}

func (s *AddressService) DeleteAddress(addressID uuid.UUID) error {
    return s.repo.Delete(&models.Address{}, addressID)
}

func (s *AddressService) GetAddressByID(addressID uuid.UUID) (*models.Address, error) {
    var address models.Address
    if err := s.repo.FindByID(&address, addressID); err != nil {
        return nil, err
    }
    return &address, nil
}