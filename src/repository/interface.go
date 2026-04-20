package repository

import "gorm.io/gorm"

type PgSQLRepository interface {
	Insert(req interface{})error
	FindOneWhere(obj interface{},query string,args ...interface{})error
	FindAllWhere(obj interface{},query string,args ...interface{})error
	UpdateByFields(obj interface{}, id interface{}, fields map[string]interface{}) error
	Delete(obj interface{}, id interface{}) error
	DeleteWhere(obj interface{}, query string, args ...interface{}) error
	FindByID(obj interface{}, id interface{}) error
	FindAll(obj interface{}) error
	Count(model interface{}, count *int64) error
	GetDB() *gorm.DB
}
