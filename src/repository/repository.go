package repository

import (
    "os"
    "gorm.io/gorm"
)

type Repository struct {
    DB *gorm.DB
}

func SetUpRepo(db *gorm.DB) *Repository {
    return &Repository{DB: db}
}

// ⚠️ FIX: Return PgSQLRepository interface, not *Repository
func (r *Repository) BeginTransaction() PgSQLRepository {
    return &Repository{DB: r.DB.Begin()}
}

func (r *Repository) Commit() error {
    return r.DB.Commit().Error
}

func (r *Repository) Rollback() error {
    return r.DB.Rollback().Error
}

// Other methods remain the same...
func (r *Repository) Insert(req interface{}) error {
    query := r.DB
    if isDevelopment() {
        query = query.Debug()
    }
    return query.Create(req).Error
}

func (r *Repository) UpdateByFields(obj interface{}, id interface{}, fields map[string]interface{}) error {
    query := r.DB
    if isDevelopment() {
        query = query.Debug()
    }
    return query.Model(obj).Where("id = ?", id).Updates(fields).Error
}

func (r *Repository) UpdateByFieldsWhere(obj interface{}, fields map[string]interface{}, query string, args ...interface{}) error {
    db := r.DB
    if isDevelopment() {
        db = db.Debug()
    }
    return db.Model(obj).Where(query, args...).Updates(fields).Error
}

func (r *Repository) Delete(obj interface{}, id interface{}) error {
    query := r.DB
    if isDevelopment() {
        query = query.Debug()
    }
    return query.Where("id = ?", id).Delete(obj).Error
}

func (r *Repository) DeleteWhere(obj interface{}, query string, args ...interface{}) error {
    db := r.DB
    if isDevelopment() {
        db = db.Debug()
    }
    return db.Where(query, args...).Delete(obj).Error
}

func (r *Repository) FindByID(obj interface{}, id interface{}) error {
    return r.DB.First(obj, "id = ?", id).Error
}

func (r *Repository) FindAll(obj interface{}) error {
    return r.DB.Find(obj).Error
}

func (r *Repository) FindOneWhere(obj interface{}, query string, args ...interface{}) error {
    return r.DB.Where(query, args...).First(obj).Error
}

func (r *Repository) FindAllWhere(obj interface{}, query string, args ...interface{}) error {
    return r.DB.Where(query, args...).Find(obj).Error
}

func (r *Repository) Count(model interface{}, count *int64) error {
    return r.DB.Model(model).Count(count).Error
}

func (r *Repository) GetDB() *gorm.DB {
    return r.DB
}

func isDevelopment() bool {
    env := os.Getenv("GO_ENV")
    return env != "production"
}