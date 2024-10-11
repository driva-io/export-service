package database

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5"
)

type PgxPostgresAdapter struct {
	Conn *pgx.Conn
}

var _ DatabaseAdapter = (*PgxPostgresAdapter)(nil)

func NewPgxPostgresAdapter(host, port, user, password, dbname string) (*PgxPostgresAdapter, error) {
	escapedUser := url.QueryEscape(user)
	escapedPassword := url.QueryEscape(password)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", escapedUser, escapedPassword, host, port, dbname)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return &PgxPostgresAdapter{Conn: conn}, nil
}

func (g *PgxPostgresAdapter) Query(query string, values []any) ([]map[string]interface{}, error) {
	rows, err := g.Conn.Query(context.Background(), query, values...)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}
	defer rows.Close()
	var results []map[string]any
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("unable to read row values: %w", err)
		}

		columns := rows.FieldDescriptions()
		rowMap := make(map[string]any)
		for i, col := range columns {
			rowMap[string(col.Name)] = values[i]
		}
		results = append(results, rowMap)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading rows: %w", rows.Err())
	}

	return results, nil
}

func (g *PgxPostgresAdapter) Close() error {
	err := g.Conn.Close(context.Background())
	if err != nil {
		log.Printf("Error closing connection.")
		return err
	}
	return nil
}
