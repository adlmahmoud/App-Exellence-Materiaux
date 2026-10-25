package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDataDirUsesExistingDatabase(t *testing.T) {
	projectDir := t.TempDir()
	legacyDBPath := filepath.Join(projectDir, "exellence.db")
	if err := os.WriteFile(legacyDBPath, []byte("existing database"), 0600); err != nil {
		t.Fatal(err)
	}

	dataDir, err := resolveDataDir(projectDir, filepath.Join(t.TempDir(), "user-config"))
	if err != nil {
		t.Fatalf("resolveDataDir returned an error: %v", err)
	}
	if dataDir != projectDir {
		t.Fatalf("resolveDataDir = %q, want existing project directory %q", dataDir, projectDir)
	}
}

func TestResolveDataDirCreatesUserDataDirectory(t *testing.T) {
	projectDir := t.TempDir()
	userConfigDir := filepath.Join(t.TempDir(), "user-config")

	dataDir, err := resolveDataDir(projectDir, userConfigDir)
	if err != nil {
		t.Fatalf("resolveDataDir returned an error: %v", err)
	}

	want := filepath.Join(userConfigDir, "Gestion de Stock")
	if dataDir != want {
		t.Fatalf("resolveDataDir = %q, want %q", dataDir, want)
	}
	if info, err := os.Stat(dataDir); err != nil || !info.IsDir() {
		t.Fatalf("user data directory was not created: %v", err)
	}
}

func TestRolePermissions(t *testing.T) {
	app := &App{}
	if err := app.requireRoles("admin"); err == nil {
		t.Fatal("unauthenticated session was granted access")
	}

	app.currentRole = "user"
	if err := app.requireRoles("admin"); err == nil {
		t.Fatal("user role was granted admin access")
	}
	if err := app.requireRoles("admin", "user"); err != nil {
		t.Fatalf("user role was denied cashier access: %v", err)
	}

	app.currentRole = "admin"
	if err := app.requireRoles("admin"); err != nil {
		t.Fatalf("admin role was denied admin access: %v", err)
	}
}
