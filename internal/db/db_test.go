package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAndPragmas(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer database.Close()

	var journalMode string
	if err := database.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("query journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode wal, got %s", journalMode)
	}

	var foreignKeys int
	if err := database.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Errorf("expected foreign_keys 1, got %d", foreignKeys)
	}

	var synchronous int
	if err := database.QueryRow("PRAGMA synchronous;").Scan(&synchronous); err != nil {
		t.Fatalf("query synchronous: %v", err)
	}
	if synchronous != 1 {
		t.Errorf("expected synchronous 1 (NORMAL), got %d", synchronous)
	}

	var version int
	if err := database.QueryRow("SELECT version FROM schema_migrations WHERE version = 1").Scan(&version); err != nil {
		t.Fatalf("query migration 1: %v", err)
	}
	if version != 1 {
		t.Errorf("expected migration version 1, got %d", version)
	}

	var authMigration int
	if err := database.QueryRow("SELECT version FROM schema_migrations WHERE version = 2").Scan(&authMigration); err != nil {
		t.Fatalf("query migration 2: %v", err)
	}
	if authMigration != 2 {
		t.Errorf("expected migration version 2, got %d", authMigration)
	}

	var authTable string
	if err := database.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'auth_credentials'").Scan(&authTable); err != nil {
		t.Fatalf("query auth_credentials table: %v", err)
	}

	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM parametres_entreprise").Scan(&count)
	if err != nil {
		t.Fatalf("query parametres_entreprise count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 row in parametres_entreprise, got %d", count)
	}

	backupPath := filepath.Join(tempDir, "backup.db")
	if err := database.Backup(backupPath); err != nil {
		t.Fatalf("Backup failed: %v", err)
	}
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file not created: %v", err)
	}
}
