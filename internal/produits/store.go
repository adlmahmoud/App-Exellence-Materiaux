package produits

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(database *sql.DB) *Store {
	return &Store{db: database}
}

func (s *Store) Create(input NouveauProduit) (Produit, error) {
	var p Produit

	if strings.TrimSpace(input.CodeBarres) == "" {
		return p, fmt.Errorf("le code-barres est obligatoire")
	}
	if strings.TrimSpace(input.Designation) == "" {
		return p, fmt.Errorf("la désignation est obligatoire")
	}

	prixTTC := input.PrixVenteTTC
	if prixTTC == 0 && input.PrixVenteHT > 0 {
		prixTTC = input.PrixVenteHT + (input.PrixVenteHT*input.TauxTVA)/10000
	}

	unite := input.Unite
	if strings.TrimSpace(unite) == "" {
		unite = "piece"
	}

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)
	id := uuid.NewString()

	tx, err := s.db.Begin()
	if err != nil {
		return p, fmt.Errorf("begin tx: %w", err)
	}

	query := `INSERT INTO produits (
		id, code_barres, reference, designation, description,
		prix_achat_ht, prix_vente_ht, prix_vente_ttc, taux_tva,
		stock_actuel, stock_minimum, unite, actif, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`

	_, err = tx.Exec(query,
		id,
		strings.TrimSpace(input.CodeBarres),
		strings.TrimSpace(input.Reference),
		strings.TrimSpace(input.Designation),
		strings.TrimSpace(input.Description),
		input.PrixAchatHT,
		input.PrixVenteHT,
		prixTTC,
		input.TauxTVA,
		input.StockInitial,
		input.StockMinimum,
		unite,
		nowStr,
		nowStr,
	)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "UNIQUE") {
			return p, fmt.Errorf("ce code-barres existe déjà")
		}
		return p, fmt.Errorf("insert produit: %w", err)
	}

	if input.StockInitial > 0 {
		mouvementID := uuid.NewString()
		mvtQuery := `INSERT INTO mouvements_stock (
			id, produit_id, type_mouvement, quantite, stock_avant, stock_apres,
			prix_unitaire_ht, reference_document, motif, created_at, updated_at
		) VALUES (?, ?, 'entree', ?, 0, ?, ?, 'INIT', 'Stock initial', ?, ?)`

		_, err = tx.Exec(mvtQuery,
			mouvementID,
			id,
			input.StockInitial,
			input.StockInitial,
			input.PrixAchatHT,
			nowStr,
			nowStr,
		)
		if err != nil {
			tx.Rollback()
			return p, fmt.Errorf("insert mouvement initial: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return p, fmt.Errorf("commit tx: %w", err)
	}

	return s.GetByID(id)
}

func (s *Store) GetByID(id string) (Produit, error) {
	var p Produit
	var createdStr, updatedStr string

	query := `SELECT id, code_barres, reference, designation, description,
		prix_achat_ht, prix_vente_ht, prix_vente_ttc, taux_tva,
		stock_actuel, stock_minimum, unite, actif, created_at, updated_at
		FROM produits WHERE id = ? LIMIT 1`

	err := s.db.QueryRow(query, id).Scan(
		&p.ID, &p.CodeBarres, &p.Reference, &p.Designation, &p.Description,
		&p.PrixAchatHT, &p.PrixVenteHT, &p.PrixVenteTTC, &p.TauxTVA,
		&p.StockActuel, &p.StockMinimum, &p.Unite, &p.Actif,
		&createdStr, &updatedStr,
	)
	if err != nil {
		return p, fmt.Errorf("produit introuvable: %w", err)
	}

	p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	return p, nil
}

func (s *Store) GetByCodeBarres(codeBarres string) (Produit, error) {
	var p Produit
	var createdStr, updatedStr string

	query := `SELECT id, code_barres, reference, designation, description,
		prix_achat_ht, prix_vente_ht, prix_vente_ttc, taux_tva,
		stock_actuel, stock_minimum, unite, actif, created_at, updated_at
		FROM produits WHERE code_barres = ? AND actif = 1 LIMIT 1`

	err := s.db.QueryRow(query, strings.TrimSpace(codeBarres)).Scan(
		&p.ID, &p.CodeBarres, &p.Reference, &p.Designation, &p.Description,
		&p.PrixAchatHT, &p.PrixVenteHT, &p.PrixVenteTTC, &p.TauxTVA,
		&p.StockActuel, &p.StockMinimum, &p.Unite, &p.Actif,
		&createdStr, &updatedStr,
	)
	if err != nil {
		return p, fmt.Errorf("produit avec ce code-barres introuvable")
	}

	p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	return p, nil
}

func (s *Store) Update(input ModificationProduit) (Produit, error) {
	var p Produit

	if strings.TrimSpace(input.ID) == "" {
		return p, fmt.Errorf("identifiant manquant")
	}
	if strings.TrimSpace(input.CodeBarres) == "" {
		return p, fmt.Errorf("le code-barres est obligatoire")
	}
	if strings.TrimSpace(input.Designation) == "" {
		return p, fmt.Errorf("la désignation est obligatoire")
	}

	prixTTC := input.PrixVenteTTC
	if prixTTC == 0 && input.PrixVenteHT > 0 {
		prixTTC = input.PrixVenteHT + (input.PrixVenteHT*input.TauxTVA)/10000
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)

	query := `UPDATE produits SET
		code_barres = ?,
		reference = ?,
		designation = ?,
		description = ?,
		prix_achat_ht = ?,
		prix_vente_ht = ?,
		prix_vente_ttc = ?,
		taux_tva = ?,
		stock_minimum = ?,
		unite = ?,
		actif = ?,
		updated_at = ?
		WHERE id = ?`

	res, err := s.db.Exec(query,
		strings.TrimSpace(input.CodeBarres),
		strings.TrimSpace(input.Reference),
		strings.TrimSpace(input.Designation),
		strings.TrimSpace(input.Description),
		input.PrixAchatHT,
		input.PrixVenteHT,
		prixTTC,
		input.TauxTVA,
		input.StockMinimum,
		input.Unite,
		input.Actif,
		nowStr,
		input.ID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return p, fmt.Errorf("ce code-barres existe déjà")
		}
		return p, fmt.Errorf("update produit: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return p, fmt.Errorf("produit introuvable")
	}

	return s.GetByID(input.ID)
}

func (s *Store) Delete(id string) error {
	query := `UPDATE produits SET actif = 0, updated_at = datetime('now') WHERE id = ?`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("suppression produit: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("produit introuvable")
	}
	return nil
}

func (s *Store) List(filter FiltreProduits) (PageProduits, error) {
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

	if filter.Actif != nil {
		whereClauses = append(whereClauses, "p.actif = ?")
		args = append(args, *filter.Actif)
	} else {
		whereClauses = append(whereClauses, "p.actif = 1")
	}

	if filter.AlerteStock {
		whereClauses = append(whereClauses, "p.stock_actuel <= p.stock_minimum")
	}

	recherche := strings.TrimSpace(filter.Recherche)
	if recherche != "" {
		words := strings.Fields(recherche)
		var ftsTokens []string
		for _, w := range words {
			var cleaned strings.Builder
			for _, r := range w {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					cleaned.WriteRune(r)
				}
			}
			sClean := cleaned.String()
			if sClean != "" {
				ftsTokens = append(ftsTokens, sClean+"*")
			}
		}

		ftsMatch := ""
		if len(ftsTokens) > 0 {
			ftsMatch = strings.Join(ftsTokens, " AND ")
		}

		var tokenClauses []string
		var tokenArgs []interface{}

		for _, w := range words {
			pat := "%" + w + "%"
			tokenClauses = append(tokenClauses, "(p.code_barres LIKE ? OR p.reference LIKE ? OR p.designation LIKE ? OR p.description LIKE ?)")
			tokenArgs = append(tokenArgs, pat, pat, pat, pat)
		}

		likePart := strings.Join(tokenClauses, " AND ")

		if ftsMatch != "" {
			combinedClause := fmt.Sprintf("((%s) OR (p.rowid IN (SELECT rowid FROM produits_fts WHERE produits_fts MATCH ?)))", likePart)
			whereClauses = append(whereClauses, combinedClause)
			args = append(args, tokenArgs...)
			args = append(args, ftsMatch)
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("(%s)", likePart))
			args = append(args, tokenArgs...)
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM produits p %s", whereSQL)
	var total int64
	if err := s.db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return PageProduits{}, fmt.Errorf("count produits: %w", err)
	}

	orderSQL := "ORDER BY p.designation ASC"
	switch filter.TriChamp {
	case "codeBarres":
		if filter.TriAscendant {
			orderSQL = "ORDER BY p.code_barres ASC"
		} else {
			orderSQL = "ORDER BY p.code_barres DESC"
		}
	case "prixVenteHT":
		if filter.TriAscendant {
			orderSQL = "ORDER BY p.prix_vente_ht ASC"
		} else {
			orderSQL = "ORDER BY p.prix_vente_ht DESC"
		}
	case "stock":
		if filter.TriAscendant {
			orderSQL = "ORDER BY p.stock_actuel ASC"
		} else {
			orderSQL = "ORDER BY p.stock_actuel DESC"
		}
	case "createdAt":
		if filter.TriAscendant {
			orderSQL = "ORDER BY p.created_at ASC"
		} else {
			orderSQL = "ORDER BY p.created_at DESC"
		}
	}

	querySQL := fmt.Sprintf(`SELECT p.id, p.code_barres, p.reference, p.designation, p.description,
		p.prix_achat_ht, p.prix_vente_ht, p.prix_vente_ttc, p.taux_tva,
		p.stock_actuel, p.stock_minimum, p.unite, p.actif, p.created_at, p.updated_at
		FROM produits p %s %s LIMIT ? OFFSET ?`, whereSQL, orderSQL)

	queryArgs := append(args, limit, offset)
	rows, err := s.db.Query(querySQL, queryArgs...)
	if err != nil {
		return PageProduits{}, fmt.Errorf("query produits: %w", err)
	}
	defer rows.Close()

	var resultats []Produit
	for rows.Next() {
		var p Produit
		var createdStr, updatedStr string
		err := rows.Scan(
			&p.ID, &p.CodeBarres, &p.Reference, &p.Designation, &p.Description,
			&p.PrixAchatHT, &p.PrixVenteHT, &p.PrixVenteTTC, &p.TauxTVA,
			&p.StockActuel, &p.StockMinimum, &p.Unite, &p.Actif,
			&createdStr, &updatedStr,
		)
		if err != nil {
			return PageProduits{}, fmt.Errorf("scan produit: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		resultats = append(resultats, p)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PageProduits{
		Produits:   resultats,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
