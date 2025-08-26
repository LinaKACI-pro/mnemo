package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/LinaKACI-pro/mnemo/internal/store"
)

const schemaEntries = `
	CREATE TABLE IF NOT EXISTS entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT NOT NULL
	);`

const queryInsert = `INSERT INTO entries (value) VALUES (?)`

const querySearch = `SELECT id, value FROM entries WHERE value LIKE ?`

type SQLiteStore struct {
	db *sql.DB
}

// New opens a database with the given driver and DSN, and initializes schema.
func New(driver, dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Create schema
	if _, err := db.Exec(schemaEntries); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Insert(ctx context.Context, e store.Entry) error {
	_, err := s.db.ExecContext(ctx, queryInsert, e.Value)
	return err
}

func (s *SQLiteStore) Search(ctx context.Context, query string) ([]store.Entry, error) {
	rows, err := s.db.QueryContext(ctx, querySearch, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []store.Entry
	for rows.Next() {
		var e store.Entry
		if err := rows.Scan(&e.ID, &e.Value); err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
