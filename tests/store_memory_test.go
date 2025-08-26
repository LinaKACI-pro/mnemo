package tests

import (
	"context"
	"testing"

	"github.com/LinaKACI-pro/mnemo/internal/store"
	"github.com/LinaKACI-pro/mnemo/store/sqlite"
)

// Smoke test using in-memory SQLite database
func TestSQLiteStore_Memory(t *testing.T) {
	// :memory: => ephemeral DB in RAM
	db, err := sqlite.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to init sqlite store: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Insert
	if err := db.Insert(ctx, store.Entry{Value: "hello"}); err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	// Search
	results, err := db.Search(ctx, "hel")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) == 0 || results[0].Value != "hello" {
		t.Fatalf("unexpected search results: %+v", results)
	}
}
