package stock

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSupprimerMouvementOnlyDeletesHistory(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE produits (id TEXT PRIMARY KEY, stock_actuel INTEGER NOT NULL);
		CREATE TABLE mouvements_stock (id TEXT PRIMARY KEY, produit_id TEXT NOT NULL);
		INSERT INTO produits (id, stock_actuel) VALUES ('produit-1', 12);
		INSERT INTO mouvements_stock (id, produit_id) VALUES ('mouvement-1', 'produit-1');
	`)
	if err != nil {
		t.Fatal(err)
	}

	store := NewStore(db)
	if err := store.SupprimerMouvement("mouvement-1"); err != nil {
		t.Fatalf("SupprimerMouvement returned error: %v", err)
	}

	var movementCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM mouvements_stock").Scan(&movementCount); err != nil {
		t.Fatal(err)
	}
	if movementCount != 0 {
		t.Errorf("expected movement history to be deleted, got %d rows", movementCount)
	}

	var stock int
	if err := db.QueryRow("SELECT stock_actuel FROM produits WHERE id = 'produit-1'").Scan(&stock); err != nil {
		t.Fatal(err)
	}
	if stock != 12 {
		t.Errorf("expected current stock to remain unchanged at 12, got %d", stock)
	}

	if err := store.SupprimerMouvement("mouvement-1"); err == nil {
		t.Error("expected deleting a missing movement to return an error")
	}
}
