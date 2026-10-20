package db

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	*sql.DB
	path string
}

func Open(dbPath string) (*DB, error) {
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	database.SetMaxOpenConns(1)

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, pragma := range pragmas {
		if _, err := database.Exec(pragma); err != nil {
			database.Close()
			return nil, fmt.Errorf("exec pragma %s: %w", pragma, err)
		}
	}

	dbInstance := &DB{
		DB:   database,
		path: dbPath,
	}

	if err := dbInstance.Migrate(); err != nil {
		database.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	dbInstance.importKbisIfPresent(filepath.Dir(dbPath))

	return dbInstance, nil
}

func (d *DB) Migrate() error {
	createMigrationsTable := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT NOT NULL
	);`

	if _, err := d.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		var version int
		_, err := fmt.Sscanf(filename, "%d_", &version)
		if err != nil {
			continue
		}

		var count int
		err = d.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration version %d: %w", version, err)
		}

		if count > 0 {
			continue
		}

		content, err := migrationsFS.ReadFile(filepath.ToSlash(filepath.Join("migrations", filename)))
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", filename, err)
		}

		tx, err := d.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction for migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("exec migration %s: %w", filename, err)
		}

		recordMigration := `INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)`
		if _, err := tx.Exec(recordMigration, version, filename, time.Now().UTC().Format(time.RFC3339)); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", filename, err)
		}
	}

	return nil
}

func (d *DB) importKbisIfPresent(baseDir string) {
	kbisPath := filepath.Join(baseDir, "kbis.json")
	content, err := os.ReadFile(kbisPath)
	if err != nil {
		return
	}

	var data struct {
		Denomination   string `json:"denomination"`
		FormeJuridique string `json:"formeJuridique"`
		Capital        int64  `json:"capital"`
		Adresse        string `json:"adresse"`
		CodePostal     string `json:"codePostal"`
		Ville          string `json:"ville"`
		Siren          string `json:"siren"`
		Siret          string `json:"siret"`
		RCS            string `json:"rcs"`
		NumeroTVA      string `json:"numeroTva"`
		Telephone      string `json:"telephone"`
		Email          string `json:"email"`
		SiteWeb        string `json:"siteWeb"`
		Devise         string `json:"devise"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return
	}

	updateQuery := `
	UPDATE parametres_entreprise SET
		denomination = CASE WHEN ? != '' THEN ? ELSE denomination END,
		forme_juridique = CASE WHEN ? != '' THEN ? ELSE forme_juridique END,
		capital = CASE WHEN ? > 0 THEN ? ELSE capital END,
		adresse = CASE WHEN ? != '' THEN ? ELSE adresse END,
		code_postal = CASE WHEN ? != '' THEN ? ELSE code_postal END,
		ville = CASE WHEN ? != '' THEN ? ELSE ville END,
		siren = CASE WHEN ? != '' THEN ? ELSE siren END,
		siret = CASE WHEN ? != '' THEN ? ELSE siret END,
		rcs = CASE WHEN ? != '' THEN ? ELSE rcs END,
		numero_tva = CASE WHEN ? != '' THEN ? ELSE numero_tva END,
		telephone = CASE WHEN ? != '' THEN ? ELSE telephone END,
		email = CASE WHEN ? != '' THEN ? ELSE email END,
		site_web = CASE WHEN ? != '' THEN ? ELSE site_web END,
		devise = CASE WHEN ? != '' THEN ? ELSE devise END,
		updated_at = datetime('now')
	WHERE id = '01DEFAULTCONFIG00000000000';
	`
	_, _ = d.Exec(updateQuery,
		data.Denomination, data.Denomination,
		data.FormeJuridique, data.FormeJuridique,
		data.Capital, data.Capital,
		data.Adresse, data.Adresse,
		data.CodePostal, data.CodePostal,
		data.Ville, data.Ville,
		data.Siren, data.Siren,
		data.Siret, data.Siret,
		data.RCS, data.RCS,
		data.NumeroTVA, data.NumeroTVA,
		data.Telephone, data.Telephone,
		data.Email, data.Email,
		data.SiteWeb, data.SiteWeb,
		data.Devise, data.Devise,
	)
}

func (d *DB) Backup(destinationPath string) error {
	cleanedPath := strings.ReplaceAll(destinationPath, "'", "''")
	cleanedPath = filepath.ToSlash(cleanedPath)
	query := fmt.Sprintf("VACUUM INTO '%s';", cleanedPath)
	if _, err := d.Exec(query); err != nil {
		return fmt.Errorf("backup database: %w", err)
	}
	return nil
}
