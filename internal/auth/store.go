package auth

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Store struct {
	db *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{db: database}
}

func (s *Store) IsConfigured() (bool, error) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM auth_credentials").Scan(&count); err != nil {
		return false, fmt.Errorf("check access credentials: %w", err)
	}
	if count != 0 && count != 2 {
		return false, fmt.Errorf("incomplete access credentials configuration")
	}
	return count == 2, nil
}

func (s *Store) SetupPasswords(adminPassword, userPassword string) error {
	if err := validatePassword(adminPassword); err != nil {
		return fmt.Errorf("admin password: %w", err)
	}
	if err := validatePassword(userPassword); err != nil {
		return fmt.Errorf("user password: %w", err)
	}
	if adminPassword == userPassword {
		return fmt.Errorf("admin and user passwords must be different")
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	userHash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash user password: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin credential setup: %w", err)
	}
	defer tx.Rollback()

	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM auth_credentials").Scan(&count); err != nil {
		return fmt.Errorf("check existing credentials: %w", err)
	}
	if count != 0 {
		return fmt.Errorf("access passwords have already been configured")
	}

	now := time.Now().UTC().Format(time.RFC3339)
	for _, credential := range []struct {
		role string
		hash []byte
	}{
		{RoleAdmin, adminHash},
		{RoleUser, userHash},
	} {
		if _, err := tx.Exec(
			"INSERT INTO auth_credentials (role, password_hash, updated_at) VALUES (?, ?, ?)",
			credential.role, string(credential.hash), now,
		); err != nil {
			return fmt.Errorf("save %s password: %w", credential.role, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit credential setup: %w", err)
	}
	return nil
}

func (s *Store) Authenticate(role, password string) error {
	if !validRole(role) {
		return fmt.Errorf("invalid role or password")
	}

	var passwordHash string
	err := s.db.QueryRow("SELECT password_hash FROM auth_credentials WHERE role = ?", role).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid role or password")
		}
		return fmt.Errorf("load %s credentials: %w", role, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return fmt.Errorf("invalid role or password")
	}
	return nil
}

func (s *Store) ChangePasswords(currentAdminPassword, adminPassword, userPassword string) error {
	if err := s.Authenticate(RoleAdmin, currentAdminPassword); err != nil {
		return fmt.Errorf("current admin password is incorrect")
	}
	if err := validatePassword(adminPassword); err != nil {
		return fmt.Errorf("admin password: %w", err)
	}
	if err := validatePassword(userPassword); err != nil {
		return fmt.Errorf("user password: %w", err)
	}
	if adminPassword == userPassword {
		return fmt.Errorf("admin and user passwords must be different")
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	userHash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash user password: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin password change: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	for _, credential := range []struct {
		role string
		hash []byte
	}{
		{RoleAdmin, adminHash},
		{RoleUser, userHash},
	} {
		if _, err := tx.Exec(
			"UPDATE auth_credentials SET password_hash = ?, updated_at = ? WHERE role = ?",
			string(credential.hash), now, credential.role,
		); err != nil {
			return fmt.Errorf("update %s password: %w", credential.role, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit password change: %w", err)
	}
	return nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if strings.TrimSpace(password) == "" || length < 8 {
		return fmt.Errorf("password must contain at least 8 characters")
	}
	if len([]byte(password)) > 72 {
		return fmt.Errorf("password must not exceed 72 bytes")
	}
	return nil
}

func validRole(role string) bool {
	return role == RoleAdmin || role == RoleUser
}
