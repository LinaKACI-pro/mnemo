package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/LinaKACI-pro/mnemo/internal/store"
)

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS documents (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        body TEXT,
        occurrence INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`,

	`CREATE VIRTUAL TABLE IF NOT EXISTS documents_fts USING fts5(
        title, body,
        content='documents',
        content_rowid='id'
    );`,

	`CREATE TRIGGER IF NOT EXISTS docs_after_insert 
    AFTER INSERT ON documents BEGIN
        INSERT INTO documents_fts(rowid, title, body)
        VALUES (new.id, new.title, new.body);
    END;`,

	`CREATE TRIGGER IF NOT EXISTS docs_after_delete 
    AFTER DELETE ON documents BEGIN
        INSERT INTO documents_fts(documents_fts, rowid, title, body)
        VALUES('delete', old.id, old.title, old.body);
    END;`,
}

const queryInsert = `
INSERT INTO documents (title, body, occurrence)
VALUES ($1, $2, $3)
`

// queryGlobalSearch performs a full-text search on the FTS5 index.
// The bm25() function computes a relevance score (the lower the score,
// the more relevant the document). Results are ordered by this score.
const queryGlobalSearch = `
SELECT docs.id, docs.title, docs.body, docs.occurrence, bm25(documents_fts) AS score
FROM documents docs
JOIN documents_fts docs_fts ON docs.id = docs_fts.rowid
WHERE documents_fts MATCH (?)
ORDER BY score;`

const querySearch = `
SELECT rowid, title, body FROM documents_fts WHERE documents_fts MATCH ?
`

type SQLiteStore struct {
	db *sql.DB
}

// New opens a database with the given driver and DSN, and initializes schema.
func New(driver, dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open(driver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	for _, stmt := range schemaStatements {
		if _, err := db.Exec(stmt); err != nil {
			return nil, fmt.Errorf("db.Exec(%q): %w", stmt, err)
		}
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Insert(ctx context.Context, title string, body string, occurrence int) error {
	_, err := s.db.ExecContext(ctx, queryInsert,
		title,
		body,
		occurrence,
	)

	return err
}

func (s *SQLiteStore) Search(ctx context.Context, metadata string, query string) ([]store.Documents, error) {
	rows, err := s.db.QueryContext(ctx, querySearch, metadata, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []store.Documents
	for rows.Next() {
		var d store.Documents
		if err := rows.Scan(
			&d.Uuid,
			&d.Title,
			&d.Body,
			&d.Occurrence,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

// GlobalSearch every FTS5 indexed fields
func (s *SQLiteStore) GlobalSearch(ctx context.Context, query string) ([]store.Documents, error) {
	rows, err := s.db.QueryContext(ctx, queryGlobalSearch, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []store.Documents
	for rows.Next() {
		var d store.Documents
		if err := rows.Scan(
			&d.Uuid,
			&d.Title,
			&d.Body,
			&d.Occurrence,
			&d.Score,
		); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		results = append(results, d)
	}
	return results, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
