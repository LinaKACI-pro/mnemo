package config

import (
	"os"
	"testing"
)

func TestLoad_Default(t *testing.T) {
	// Ensure no ENV variable is set
	_ = os.Unsetenv("MNEMO_DB_PATH")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DbPath != "./mnemo.db" {
		t.Errorf("expected default ./mnemo.db, got %s", cfg.DbPath)
	}
}

func TestLoad_WithEnv(t *testing.T) {
	// Set ENV variable
	testPath := "/tmp/test.db"
	_ = os.Setenv("MNEMO_DB_PATH", testPath)
	defer os.Unsetenv("MNEMO_DB_PATH")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DbPath != testPath {
		t.Errorf("expected %s from env, got %s", testPath, cfg.DbPath)
	}
}
