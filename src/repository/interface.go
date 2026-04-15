package repository

type PgSQLRepository interface {
	Insert(req interface{}) error
	FindOneWhere(obj interface{}, query string, args ...interface{}) error
	UpdateByFields(obj interface{}, id interface{}, fields map[string]interface{}) error
	Delete(obj interface{}, id interface{}) error
}
