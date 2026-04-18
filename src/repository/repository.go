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
func (r *Repository)DeleteWhere(obj interface{},query string,args ...interface{})error{
	return r.DB.Where(query,args...).Delete(obj).Error
}
func (r *Repository)FindByID(obj interface{},id ...interface{},)error{
	return r.DB.First(obj,"id = ?",id).Error
}

func (r *Repository)FindAll(obj interface{})error{
	return r.DB.Find(obj).Error
}

func (r *Repository)FindWhere(obj interface{},query string,args ...interface{})error{
	return r.DB.Where(query,args...).Find(obj).Error
}

func (r *Repository)Count(model interface{},count *int64)error{
	return r.DB.Model(model).Count(count).Error
}

func (r *Repository)GetDB()*gorm.DB{
	return r.DB
}