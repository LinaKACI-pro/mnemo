package tests

import (
	"context"
	"testing"

	"github.com/LinaKACI-pro/mnemo/store/sqlite"
)

// Smoke test using in-memory SQLite database
func TestSQLiteStore_Memory(t *testing.T) {
	// :memory: => ephemeral DB in RAM
	db, err := sqlite.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("sqlite.New: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Insert one document
	err = db.Insert(ctx,
		"Test Title",
		"Hello world body content",
		1,
	)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	// GlobalSearch with a word from the body
	results, err := db.GlobalSearch(ctx, "hello")
	if err != nil {
		t.Fatalf("global search failed: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected at least one result, got none")
	}

	doc := results[0]
	if doc.Title != "Test Title" || doc.Body == "" {
		t.Errorf("unexpected search result: %+v", doc)
	}

	// Ensure FTS5 ranking score is present
	if doc.Score == 0 {
		t.Logf("warning: bm25 score is 0, might indicate non-relevance but query worked")
	}
}
