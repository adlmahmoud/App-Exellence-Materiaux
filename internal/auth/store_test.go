package auth

import (
	"database/sql"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		database.Close()
	})
	if _, err := database.Exec(`
		CREATE TABLE auth_credentials (
			role TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)
	`); err != nil {
		t.Fatal(err)
	}
	return NewStore(database)
}

func TestSetupAuthenticateAndChangePasswords(t *testing.T) {
	store := newTestStore(t)

	configured, err := store.IsConfigured()
	if err != nil || configured {
		t.Fatalf("initial IsConfigured() = %v, %v; want false, nil", configured, err)
	}
	if err := store.SetupPasswords("admin-pass-123", "user-pass-456"); err != nil {
		t.Fatalf("SetupPasswords() error = %v", err)
	}
	var storedHash string
	if err := store.db.QueryRow("SELECT password_hash FROM auth_credentials WHERE role = ?", RoleAdmin).Scan(&storedHash); err != nil {
		t.Fatal(err)
	}
	if storedHash == "admin-pass-123" || !strings.HasPrefix(storedHash, "$2a$") {
		t.Fatal("admin password was not stored as a bcrypt hash")
	}
	configured, err = store.IsConfigured()
	if err != nil || !configured {
		t.Fatalf("configured IsConfigured() = %v, %v; want true, nil", configured, err)
	}
	if err := store.SetupPasswords("other-admin-123", "other-user-456"); err == nil {
		t.Fatal("SetupPasswords() succeeded after credentials were configured")
	}
	if err := store.Authenticate(RoleAdmin, "wrong-password"); err == nil {
		t.Fatal("Authenticate() accepted an incorrect password")
	}
	if err := store.Authenticate(RoleAdmin, "admin-pass-123"); err != nil {
		t.Fatalf("Authenticate() rejected the admin password: %v", err)
	}
	if err := store.ChangePasswords("wrong-password", "admin-new-123", "user-new-456"); err == nil {
		t.Fatal("ChangePasswords() accepted an incorrect current password")
	}
	if err := store.ChangePasswords("admin-pass-123", "admin-new-123", "user-new-456"); err != nil {
		t.Fatalf("ChangePasswords() error = %v", err)
	}
	if err := store.Authenticate(RoleAdmin, "admin-new-123"); err != nil {
		t.Fatalf("Authenticate() rejected the changed admin password: %v", err)
	}
	if err := store.Authenticate(RoleUser, "user-new-456"); err != nil {
		t.Fatalf("Authenticate() rejected the changed user password: %v", err)
	}
}

func TestSetupPasswordsValidatesPasswords(t *testing.T) {
	tests := []struct {
		name          string
		adminPassword string
		userPassword  string
	}{
		{"too short", "short", "user-pass-456"},
		{"same passwords", "same-pass-123", "same-pass-123"},
		{"too long for bcrypt", strings.Repeat("a", 73), "user-pass-456"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := newTestStore(t).SetupPasswords(test.adminPassword, test.userPassword); err == nil {
				t.Fatal("SetupPasswords() accepted invalid passwords")
			}
		})
	}
}
