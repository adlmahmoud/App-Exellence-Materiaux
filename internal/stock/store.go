package stock

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{db: database}
}

func (s *Store) CreerMouvement(input SaisieMouvement) (MouvementStock, error) {
	var m MouvementStock

	if strings.TrimSpace(input.ProduitID) == "" {
		return m, fmt.Errorf("produit non spécifié")
	}
	if input.Quantite <= 0 {
		return m, fmt.Errorf("la quantité doit être supérieure à zéro")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return m, fmt.Errorf("begin tx: %w", err)
	}

	var stockActuel int64
	err = tx.QueryRow("SELECT stock_actuel FROM produits WHERE id = ? FOR UPDATE", input.ProduitID).Scan(&stockActuel)
	if err != nil {
		err = tx.QueryRow("SELECT stock_actuel FROM produits WHERE id = ?", input.ProduitID).Scan(&stockActuel)
		if err != nil {
			tx.Rollback()
			return m, fmt.Errorf("produit introuvable: %w", err)
		}
	}

	var nouveauStock int64
	switch input.TypeMouvement {
	case MouvementEntree, MouvementRetour, MouvementReception:
		nouveauStock = stockActuel + input.Quantite
	case MouvementSortie, MouvementVente:
		nouveauStock = stockActuel - input.Quantite
	default:
		tx.Rollback()
		return m, fmt.Errorf("type de mouvement invalide: %s", input.TypeMouvement)
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	_, err = tx.Exec("UPDATE produits SET stock_actuel = ?, updated_at = ? WHERE id = ?", nouveauStock, nowStr, input.ProduitID)
	if err != nil {
		tx.Rollback()
		return m, fmt.Errorf("update stock: %w", err)
	}

	query := `INSERT INTO mouvements_stock (
		id, produit_id, type_mouvement, quantite, stock_avant, stock_apres,
		prix_unitaire_ht, reference_document, motif, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(query,
		id,
		input.ProduitID,
		input.TypeMouvement,
		input.Quantite,
		stockActuel,
		nouveauStock,
		input.PrixUnitaireHT,
		strings.TrimSpace(input.ReferenceDocument),
		strings.TrimSpace(input.Motif),
		nowStr,
		nowStr,
	)
	if err != nil {
		tx.Rollback()
		return m, fmt.Errorf("insert mouvement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return m, fmt.Errorf("commit tx: %w", err)
	}

	return s.GetByID(id)
}

func (s *Store) AjusterStock(input SaisieAjustement) (MouvementStock, error) {
	var m MouvementStock

	if strings.TrimSpace(input.ProduitID) == "" {
		return m, fmt.Errorf("produit non spécifié")
	}
	if input.NouveauStock < 0 {
		return m, fmt.Errorf("le nouveau stock ne peut pas être négatif")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return m, fmt.Errorf("begin tx: %w", err)
	}

	var stockActuel int64
	err = tx.QueryRow("SELECT stock_actuel FROM produits WHERE id = ?", input.ProduitID).Scan(&stockActuel)
	if err != nil {
		tx.Rollback()
		return m, fmt.Errorf("produit introuvable: %w", err)
	}

	delta := input.NouveauStock - stockActuel
	nowStr := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	_, err = tx.Exec("UPDATE produits SET stock_actuel = ?, updated_at = ? WHERE id = ?", input.NouveauStock, nowStr, input.ProduitID)
	if err != nil {
		tx.Rollback()
		return m, fmt.Errorf("update stock: %w", err)
	}

	motif := strings.TrimSpace(input.Motif)
	if motif == "" {
		motif = "Ajustement d'inventaire"
	}

	query := `INSERT INTO mouvements_stock (
		id, produit_id, type_mouvement, quantite, stock_avant, stock_apres,
		prix_unitaire_ht, reference_document, motif, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, 0, 'AJUST', ?, ?, ?)`

	_, err = tx.Exec(query,
		id,
		input.ProduitID,
		MouvementAjustement,
		delta,
		stockActuel,
		input.NouveauStock,
		motif,
		nowStr,
		nowStr,
	)
	if err != nil {
		tx.Rollback()
		return m, fmt.Errorf("insert mouvement ajustement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return m, fmt.Errorf("commit tx: %w", err)
	}

	return s.GetByID(id)
}

func (s *Store) GetByID(id string) (MouvementStock, error) {
	var m MouvementStock
	var createdStr, updatedStr string

	query := `SELECT id, produit_id, type_mouvement, quantite, stock_avant, stock_apres,
		prix_unitaire_ht, reference_document, motif, created_at, updated_at
		FROM mouvements_stock WHERE id = ? LIMIT 1`

	err := s.db.QueryRow(query, id).Scan(
		&m.ID, &m.ProduitID, &m.TypeMouvement, &m.Quantite, &m.StockAvant, &m.StockApres,
		&m.PrixUnitaireHT, &m.ReferenceDocument, &m.Motif, &createdStr, &updatedStr,
	)
	if err != nil {
		return m, fmt.Errorf("mouvement introuvable: %w", err)
	}

	m.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	m.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	return m, nil
}

func (s *Store) SupprimerMouvement(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("identifiant du mouvement non spécifié")
	}

	result, err := s.db.Exec("DELETE FROM mouvements_stock WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("supprimer mouvement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("vérifier la suppression du mouvement: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("mouvement introuvable")
	}
	return nil
}

func (s *Store) List(filter FiltreMouvements) (PageMouvements, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 500 {
		limit = 50
	}
	offset := (page - 1) * limit

	var whereClauses []string
	var args []interface{}

	if strings.TrimSpace(filter.ProduitID) != "" {
		whereClauses = append(whereClauses, "m.produit_id = ?")
		args = append(args, strings.TrimSpace(filter.ProduitID))
	}

	if filter.TypeMouvement != nil {
		whereClauses = append(whereClauses, "m.type_mouvement = ?")
		args = append(args, string(*filter.TypeMouvement))
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM mouvements_stock m %s", whereSQL)
	var total int64
	if err := s.db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return PageMouvements{}, fmt.Errorf("count mouvements: %w", err)
	}

	querySQL := fmt.Sprintf(`SELECT m.id, m.produit_id, m.type_mouvement, m.quantite, m.stock_avant, m.stock_apres,
		m.prix_unitaire_ht, m.reference_document, m.motif, m.created_at, m.updated_at,
		p.code_barres, p.reference, p.designation
		FROM mouvements_stock m
		JOIN produits p ON p.id = m.produit_id
		%s
		ORDER BY m.created_at DESC LIMIT ? OFFSET ?`, whereSQL)

	queryArgs := append(args, limit, offset)
	rows, err := s.db.Query(querySQL, queryArgs...)
	if err != nil {
		return PageMouvements{}, fmt.Errorf("query mouvements: %w", err)
	}
	defer rows.Close()

	var resultats []MouvementDetail
	for rows.Next() {
		var md MouvementDetail
		var createdStr, updatedStr string
		err := rows.Scan(
			&md.ID, &md.ProduitID, &md.TypeMouvement, &md.Quantite, &md.StockAvant, &md.StockApres,
			&md.PrixUnitaireHT, &md.ReferenceDocument, &md.Motif, &createdStr, &updatedStr,
			&md.CodeBarres, &md.Reference, &md.Designation,
		)
		if err != nil {
			return PageMouvements{}, fmt.Errorf("scan mouvement: %w", err)
		}
		md.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		md.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		resultats = append(resultats, md)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PageMouvements{
		Mouvements: resultats,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) ObtenirAlertes() ([]AlerteStock, error) {
	query := `SELECT id, code_barres, reference, designation, stock_actuel, stock_minimum, unite
		FROM produits
		WHERE actif = 1 AND stock_actuel <= stock_minimum
		ORDER BY (stock_actuel - stock_minimum) ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query alertes: %w", err)
	}
	defer rows.Close()

	var alertes []AlerteStock
	for rows.Next() {
		var a AlerteStock
		if err := rows.Scan(&a.ProduitID, &a.CodeBarres, &a.Reference, &a.Designation, &a.StockActuel, &a.StockMinimum, &a.Unite); err != nil {
			return nil, fmt.Errorf("scan alerte: %w", err)
		}
		alertes = append(alertes, a)
	}
	return alertes, nil
}
