package repository

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func SetUpRepo(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Insert(req interface{}) error {
	return r.DB.Create(req).Error
}

func (r *Repository) FindOneWhere(obj interface{}, query string, args ...interface{}) error {
	return r.DB.Where(query, args...).First(obj).Error
}

func (r *Repository) UpdateByFields(obj interface{}, id interface{}, fields map[string]interface{}) error {
	return r.DB.Model(obj).Where("id = ?", id).Updates(fields).Error
}

func (r *Repository) Delete(obj interface{}, id interface{}) error {
	return r.DB.Where("id = ?", id).Delete(obj).Error
}