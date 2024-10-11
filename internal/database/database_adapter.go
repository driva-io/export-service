package database

type DatabaseAdapter interface {
	Close() error
	Query(query string, values []any) ([]map[string]interface{}, error)
}
