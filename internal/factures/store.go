package factures

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

func (s *Store) ProchainNumero(tx *sql.Tx, typeDocument string, annee string) (int64, string, error) {
	cleSequence := fmt.Sprintf("%s_%s", typeDocument, annee)

	var derniereValeur int64
	err := tx.QueryRow("SELECT derniere_valeur FROM sequences WHERE nom = ?", cleSequence).Scan(&derniereValeur)
	if err != nil && err != sql.ErrNoRows {
		return 0, "", fmt.Errorf("select sequence: %w", err)
	}

	prochaineValeur := derniereValeur + 1
	nowStr := time.Now().UTC().Format(time.RFC3339)

	if err == sql.ErrNoRows {
		_, err = tx.Exec("INSERT INTO sequences (nom, derniere_valeur, updated_at) VALUES (?, ?, ?)", cleSequence, prochaineValeur, nowStr)
	} else {
		_, err = tx.Exec("UPDATE sequences SET derniere_valeur = ?, updated_at = ? WHERE nom = ?", prochaineValeur, nowStr, cleSequence)
	}
	if err != nil {
		return 0, "", fmt.Errorf("update sequence: %w", err)
	}

	var format string
	if typeDocument == "facture" {
		_ = tx.QueryRow("SELECT format_numerotation_facture FROM parametres_entreprise LIMIT 1").Scan(&format)
		if format == "" {
			format = "FAC-{YYYY}-{SEQ}"
		}
	} else {
		_ = tx.QueryRow("SELECT format_numerotation_bon FROM parametres_entreprise LIMIT 1").Scan(&format)
		if format == "" {
			format = "BON-{TYPE}-{YYYY}-{SEQ}"
		}
	}

	numero := format
	numero = strings.ReplaceAll(numero, "{YYYY}", annee)
	numero = strings.ReplaceAll(numero, "{TYPE}", strings.ToUpper(typeDocument))
	seqFormatee := fmt.Sprintf("%05d", prochaineValeur)
	numero = strings.ReplaceAll(numero, "{SEQ}", seqFormatee)

	return prochaineValeur, numero, nil
}

func (s *Store) CreerFacture(input NouvelleFacture) (DetailFacture, error) {
	var detail DetailFacture

	if len(input.Lignes) == 0 {
		return detail, fmt.Errorf("la facture doit contenir au moins une ligne")
	}

	clientNom := strings.TrimSpace(input.ClientNom)
	if clientNom == "" {
		clientNom = "Client Comptoir"
	}

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)
	annee := fmt.Sprintf("%04d", now.Year())
	factureID := uuid.NewString()

	tx, err := s.db.Begin()
	if err != nil {
		return detail, fmt.Errorf("begin tx: %w", err)
	}

	var paramDenom, paramForme, paramAdresse, paramCP, paramVille string
	var paramSiren, paramSiret, paramRCS, paramTVA, paramTel, paramEmail, paramMentions, paramConditions, paramDevise string
	var paramCapital int64

	paramQuery := `SELECT denomination, forme_juridique, capital, adresse, code_postal, ville,
		siren, siret, rcs, numero_tva, telephone, email, mentions_legales, conditions_reglement, devise
		FROM parametres_entreprise LIMIT 1`

	err = tx.QueryRow(paramQuery).Scan(
		&paramDenom, &paramForme, &paramCapital, &paramAdresse, &paramCP, &paramVille,
		&paramSiren, &paramSiret, &paramRCS, &paramTVA, &paramTel, &paramEmail,
		&paramMentions, &paramConditions, &paramDevise,
	)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("recuperation snapshot entreprise: %w", err)
	}

	seqNumero, numeroFacture, err := s.ProchainNumero(tx, "facture", annee)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("generation numero facture: %w", err)
	}

	dateFacture := input.DateFacture
	if strings.TrimSpace(dateFacture) == "" {
		dateFacture = fmt.Sprintf("%02d/%02d/%04d", now.Day(), now.Month(), now.Year())
	}
	dateEcheance := input.DateEcheance
	if strings.TrimSpace(dateEcheance) == "" {
		dateEcheance = dateFacture
	}

	var totalHT int64 = 0
	var totalTVA int64 = 0
	var totalTTC int64 = 0
	var lignesInserees []LigneFacture

	for i, l := range input.Lignes {
		if strings.TrimSpace(l.Designation) == "" {
			tx.Rollback()
			return detail, fmt.Errorf("désignation obligatoire sur la ligne %d", i+1)
		}
		if l.Quantite <= 0 {
			tx.Rollback()
			return detail, fmt.Errorf("quantité invalide sur la ligne %d", i+1)
		}

		montantHT := l.PrixUnitaireHT * l.Quantite
		montantTVA := (montantHT * l.TauxTVA) / 10000
		montantTTC := montantHT + montantTVA

		totalHT += montantHT
		totalTVA += montantTVA
		totalTTC += montantTTC

		ligneID := uuid.NewString()

		ligneQuery := `INSERT INTO lignes_facture (
			id, facture_id, produit_id, code_barres, reference, designation,
			quantite, prix_unitaire_ht, taux_tva, montant_ht, montant_tva, montant_ttc,
			ordre, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(ligneQuery,
			ligneID,
			factureID,
			l.ProduitID,
			strings.TrimSpace(l.CodeBarres),
			strings.TrimSpace(l.Reference),
			strings.TrimSpace(l.Designation),
			l.Quantite,
			l.PrixUnitaireHT,
			l.TauxTVA,
			montantHT,
			montantTVA,
			montantTTC,
			i,
			nowStr,
			nowStr,
		)
		if err != nil {
			tx.Rollback()
			return detail, fmt.Errorf("insert ligne facture %d: %w", i+1, err)
		}

		if l.ProduitID != nil && *l.ProduitID != "" {
			var stockActuel int64
			err := tx.QueryRow("SELECT stock_actuel FROM produits WHERE id = ?", *l.ProduitID).Scan(&stockActuel)
			if err == nil {
				nouveauStock := stockActuel - l.Quantite
				_, _ = tx.Exec("UPDATE produits SET stock_actuel = ?, updated_at = ? WHERE id = ?", nouveauStock, nowStr, *l.ProduitID)

				mvtID := uuid.NewString()
				mvtSQL := `INSERT INTO mouvements_stock (
					id, produit_id, type_mouvement, quantite, stock_avant, stock_apres,
					prix_unitaire_ht, reference_document, motif, created_at, updated_at
				) VALUES (?, ?, 'vente', ?, ?, ?, ?, ?, 'Vente facture', ?, ?)`

				_, _ = tx.Exec(mvtSQL, mvtID, *l.ProduitID, l.Quantite, stockActuel, nouveauStock, l.PrixUnitaireHT, numeroFacture, nowStr, nowStr)
			}
		}

		lignesInserees = append(lignesInserees, LigneFacture{
			ID:             ligneID,
			FactureID:      factureID,
			ProduitID:      l.ProduitID,
			CodeBarres:     l.CodeBarres,
			Reference:      l.Reference,
			Designation:    l.Designation,
			Quantite:       l.Quantite,
			PrixUnitaireHT: l.PrixUnitaireHT,
			TauxTVA:        l.TauxTVA,
			MontantHT:      montantHT,
			MontantTVA:     montantTVA,
			MontantTTC:     montantTTC,
			Ordre:          i,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	modeReg := input.ModeReglement
	if modeReg == "" {
		modeReg = ReglementEspeces
	}

	montantRegle := input.MontantRegle
	if montantRegle == 0 {
		montantRegle = totalTTC
	}
	resteAPayer := totalTTC - montantRegle
	if resteAPayer < 0 {
		resteAPayer = 0
	}

	statut := StatutValidee
	if resteAPayer == 0 {
		statut = StatutPayee
	}

	conditions := input.ConditionsReglement
	if strings.TrimSpace(conditions) == "" {
		conditions = paramConditions
	}

	insertFactureSQL := `INSERT INTO factures (
		id, numero, sequence_numero, statut, mode_reglement, date_facture, date_echeance, conditions_reglement,
		entreprise_denomination, entreprise_forme_juridique, entreprise_capital, entreprise_adresse,
		entreprise_code_postal, entreprise_ville, entreprise_siren, entreprise_siret, entreprise_rcs,
		entreprise_numero_tva, entreprise_telephone, entreprise_email, mentions_legales, devise,
		client_id, client_nom, client_adresse, client_siren_siret,
		total_ht, total_tva, total_ttc, montant_regle, reste_a_payer, notes,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(insertFactureSQL,
		factureID,
		numeroFacture,
		seqNumero,
		statut,
		modeReg,
		dateFacture,
		dateEcheance,
		conditions,
		paramDenom,
		paramForme,
		paramCapital,
		paramAdresse,
		paramCP,
		paramVille,
		paramSiren,
		paramSiret,
		paramRCS,
		paramTVA,
		paramTel,
		paramEmail,
		paramMentions,
		paramDevise,
		input.ClientID,
		clientNom,
		strings.TrimSpace(input.ClientAdresse),
		strings.TrimSpace(input.ClientSirenSiret),
		totalHT,
		totalTVA,
		totalTTC,
		montantRegle,
		resteAPayer,
		strings.TrimSpace(input.Notes),
		nowStr,
		nowStr,
	)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("insert facture: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return detail, fmt.Errorf("commit facture: %w", err)
	}

	return s.GetByID(factureID)
}

func (s *Store) GetByID(id string) (DetailFacture, error) {
	var detail DetailFacture
	var f Facture
	var createdStr, updatedStr string

	query := `SELECT id, numero, sequence_numero, statut, mode_reglement, date_facture, date_echeance, conditions_reglement,
		entreprise_denomination, entreprise_forme_juridique, entreprise_capital, entreprise_adresse,
		entreprise_code_postal, entreprise_ville, entreprise_siren, entreprise_siret, entreprise_rcs,
		entreprise_numero_tva, entreprise_telephone, entreprise_email, mentions_legales, devise,
		client_id, client_nom, client_adresse, client_code_postal, client_ville, client_telephone, client_email,
		client_siren_siret, client_numero_tva, total_ht, total_tva, total_ttc, montant_regle, reste_a_payer, notes,
		created_at, updated_at
		FROM factures WHERE id = ? LIMIT 1`

	err := s.db.QueryRow(query, id).Scan(
		&f.ID, &f.Numero, &f.SequenceNumero, &f.Statut, &f.ModeReglement, &f.DateFacture, &f.DateEcheance, &f.ConditionsReglement,
		&f.EntrepriseDenomination, &f.EntrepriseFormeJuridique, &f.EntrepriseCapital, &f.EntrepriseAdresse,
		&f.EntrepriseCodePostal, &f.EntrepriseVille, &f.EntrepriseSiren, &f.EntrepriseSiret, &f.EntrepriseRCS,
		&f.EntrepriseNumeroTVA, &f.EntrepriseTelephone, &f.EntrepriseEmail, &f.MentionsLegales, &f.Devise,
		&f.ClientID, &f.ClientNom, &f.ClientAdresse, &f.ClientCodePostal, &f.ClientVille, &f.ClientTelephone, &f.ClientEmail,
		&f.ClientSirenSiret, &f.ClientNumeroTVA, &f.TotalHT, &f.TotalTVA, &f.TotalTTC, &f.MontantRegle, &f.ResteAPayer, &f.Notes,
		&createdStr, &updatedStr,
	)
	if err != nil {
		return detail, fmt.Errorf("facture introuvable: %w", err)
	}

	f.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	f.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
	detail.Facture = f

	lignesQuery := `SELECT id, facture_id, produit_id, code_barres, reference, designation,
		quantite, prix_unitaire_ht, taux_tva, montant_ht, montant_tva, montant_ttc, ordre, created_at, updated_at
		FROM lignes_facture WHERE facture_id = ? ORDER BY ordre ASC`

	rows, err := s.db.Query(lignesQuery, id)
	if err != nil {
		return detail, fmt.Errorf("query lignes facture: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l LigneFacture
		var lCreatedStr, lUpdatedStr string
		err := rows.Scan(
			&l.ID, &l.FactureID, &l.ProduitID, &l.CodeBarres, &l.Reference, &l.Designation,
			&l.Quantite, &l.PrixUnitaireHT, &l.TauxTVA, &l.MontantHT, &l.MontantTVA, &l.MontantTTC, &l.Ordre,
			&lCreatedStr, &lUpdatedStr,
		)
		if err != nil {
			return detail, fmt.Errorf("scan ligne facture: %w", err)
		}
		l.CreatedAt, _ = time.Parse(time.RFC3339, lCreatedStr)
		l.UpdatedAt, _ = time.Parse(time.RFC3339, lUpdatedStr)
		detail.Lignes = append(detail.Lignes, l)
	}

	return detail, nil
}

func (s *Store) List(filter FiltreFactures) (PageFactures, error) {
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

	recherche := strings.TrimSpace(filter.Recherche)
	if recherche != "" {
		whereClauses = append(whereClauses, "(numero LIKE ? OR client_nom LIKE ?)")
		args = append(args, "%"+recherche+"%", "%"+recherche+"%")
	}

	if filter.Statut != nil {
		whereClauses = append(whereClauses, "statut = ?")
		args = append(args, string(*filter.Statut))
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM factures %s", whereSQL)
	var total int64
	if err := s.db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return PageFactures{}, fmt.Errorf("count factures: %w", err)
	}

	querySQL := fmt.Sprintf(`SELECT id, numero, sequence_numero, statut, mode_reglement, date_facture, date_echeance, conditions_reglement,
		entreprise_denomination, entreprise_forme_juridique, entreprise_capital, entreprise_adresse,
		entreprise_code_postal, entreprise_ville, entreprise_siren, entreprise_siret, entreprise_rcs,
		entreprise_numero_tva, entreprise_telephone, entreprise_email, mentions_legales, devise,
		client_id, client_nom, client_adresse, client_code_postal, client_ville, client_telephone, client_email,
		client_siren_siret, client_numero_tva, total_ht, total_tva, total_ttc, montant_regle, reste_a_payer, notes,
		created_at, updated_at
		FROM factures %s ORDER BY sequence_numero DESC LIMIT ? OFFSET ?`, whereSQL)

	queryArgs := append(args, limit, offset)
	rows, err := s.db.Query(querySQL, queryArgs...)
	if err != nil {
		return PageFactures{}, fmt.Errorf("query factures: %w", err)
	}
	defer rows.Close()

	var resultats []Facture
	for rows.Next() {
		var f Facture
		var createdStr, updatedStr string
		err := rows.Scan(
			&f.ID, &f.Numero, &f.SequenceNumero, &f.Statut, &f.ModeReglement, &f.DateFacture, &f.DateEcheance, &f.ConditionsReglement,
			&f.EntrepriseDenomination, &f.EntrepriseFormeJuridique, &f.EntrepriseCapital, &f.EntrepriseAdresse,
			&f.EntrepriseCodePostal, &f.EntrepriseVille, &f.EntrepriseSiren, &f.EntrepriseSiret, &f.EntrepriseRCS,
			&f.EntrepriseNumeroTVA, &f.EntrepriseTelephone, &f.EntrepriseEmail, &f.MentionsLegales, &f.Devise,
			&f.ClientID, &f.ClientNom, &f.ClientAdresse, &f.ClientCodePostal, &f.ClientVille, &f.ClientTelephone, &f.ClientEmail,
			&f.ClientSirenSiret, &f.ClientNumeroTVA, &f.TotalHT, &f.TotalTVA, &f.TotalTTC, &f.MontantRegle, &f.ResteAPayer, &f.Notes,
			&createdStr, &updatedStr,
		)
		if err != nil {
			return PageFactures{}, fmt.Errorf("scan facture: %w", err)
		}
		f.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		f.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		resultats = append(resultats, f)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PageFactures{
		Factures:   resultats,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Store) CreerBon(input NouveauBon) (DetailBon, error) {
	var detail DetailBon

	if len(input.Lignes) == 0 {
		return detail, fmt.Errorf("le bon doit contenir au moins une ligne")
	}

	clientNom := strings.TrimSpace(input.ClientNom)
	if clientNom == "" {
		clientNom = "Client Comptoir"
	}

	typeBon := input.TypeBon
	if typeBon == "" {
		typeBon = BonLivraison
	}

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)
	annee := fmt.Sprintf("%04d", now.Year())
	bonID := uuid.NewString()

	tx, err := s.db.Begin()
	if err != nil {
		return detail, fmt.Errorf("begin tx: %w", err)
	}

	var paramDenom, paramAdresse, paramTel, paramEmail string
	paramQuery := `SELECT denomination, adresse, telephone, email FROM parametres_entreprise LIMIT 1`
	err = tx.QueryRow(paramQuery).Scan(&paramDenom, &paramAdresse, &paramTel, &paramEmail)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("recuperation snapshot entreprise: %w", err)
	}

	seqNumero, numeroBon, err := s.ProchainNumero(tx, "bon", annee)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("generation numero bon: %w", err)
	}

	dateBon := input.DateBon
	if strings.TrimSpace(dateBon) == "" {
		dateBon = fmt.Sprintf("%02d/%02d/%04d", now.Day(), now.Month(), now.Year())
	}

	var totalHT int64 = 0
	var totalTVA int64 = 0
	var totalTTC int64 = 0
	var lignesInserees []LigneBon

	for i, l := range input.Lignes {
		if strings.TrimSpace(l.Designation) == "" {
			tx.Rollback()
			return detail, fmt.Errorf("désignation obligatoire sur la ligne %d", i+1)
		}
		if l.Quantite <= 0 {
			tx.Rollback()
			return detail, fmt.Errorf("quantité invalide sur la ligne %d", i+1)
		}

		montantHT := l.PrixUnitaireHT * l.Quantite
		montantTVA := (montantHT * l.TauxTVA) / 10000
		montantTTC := montantHT + montantTVA

		totalHT += montantHT
		totalTVA += montantTVA
		totalTTC += montantTTC

		ligneID := uuid.NewString()
		ligneQuery := `INSERT INTO lignes_bon (
			id, bon_id, produit_id, code_barres, reference, designation,
			quantite, prix_unitaire_ht, taux_tva, montant_ht, montant_tva, montant_ttc, ordre,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		_, err = tx.Exec(ligneQuery,
			ligneID, bonID, l.ProduitID, l.CodeBarres, l.Reference, l.Designation,
			l.Quantite, l.PrixUnitaireHT, l.TauxTVA, montantHT, montantTVA, montantTTC, i+1,
			nowStr, nowStr,
		)
		if err != nil {
			tx.Rollback()
			return detail, fmt.Errorf("insert ligne bon: %w", err)
		}

		lignesInserees = append(lignesInserees, LigneBon{
			ID:             ligneID,
			BonID:          bonID,
			ProduitID:      l.ProduitID,
			CodeBarres:     l.CodeBarres,
			Reference:      l.Reference,
			Designation:    l.Designation,
			Quantite:       l.Quantite,
			PrixUnitaireHT: l.PrixUnitaireHT,
			TauxTVA:        l.TauxTVA,
			MontantHT:      montantHT,
			MontantTVA:     montantTVA,
			MontantTTC:     montantTTC,
			Ordre:          i + 1,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}

	bonQuery := `INSERT INTO bons (
		id, type_bon, numero, sequence_numero, statut, date_bon, reference_externe,
		entreprise_denomination, entreprise_adresse, entreprise_telephone, entreprise_email,
		client_id, client_nom, total_ht, total_tva, total_ttc, notes, created_at, updated_at
	) VALUES (?, ?, ?, ?, 'en_cours', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(bonQuery,
		bonID, string(typeBon), numeroBon, seqNumero, dateBon, input.ReferenceExterne,
		paramDenom, paramAdresse, paramTel, paramEmail,
		input.ClientID, clientNom, totalHT, totalTVA, totalTTC, input.Notes,
		nowStr, nowStr,
	)
	if err != nil {
		tx.Rollback()
		return detail, fmt.Errorf("insert bon: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return detail, fmt.Errorf("commit tx: %w", err)
	}

	detail.Bon = Bon{
		ID:                     bonID,
		TypeBon:                typeBon,
		Numero:                 numeroBon,
		SequenceNumero:         seqNumero,
		Statut:                 "en_cours",
		DateBon:                dateBon,
		ReferenceExterne:       input.ReferenceExterne,
		EntrepriseDenomination: paramDenom,
		EntrepriseAdresse:      paramAdresse,
		EntrepriseTelephone:    paramTel,
		EntrepriseEmail:        paramEmail,
		ClientID:               input.ClientID,
		ClientNom:              clientNom,
		TotalHT:                totalHT,
		TotalTVA:               totalTVA,
		TotalTTC:               totalTTC,
		Notes:                  input.Notes,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	detail.Lignes = lignesInserees

	return detail, nil
}

func (s *Store) GetBonByID(id string) (DetailBon, error) {
	var detail DetailBon

	query := `SELECT id, type_bon, numero, sequence_numero, statut, date_bon, reference_externe,
		facture_id, entreprise_denomination, entreprise_adresse, entreprise_telephone, entreprise_email,
		client_id, client_nom, client_adresse, client_code_postal, client_ville, client_telephone,
		total_ht, total_tva, total_ttc, notes, created_at, updated_at
		FROM bons WHERE id = ? LIMIT 1`

	var createdStr, updatedStr string
	var typeBonStr string
	err := s.db.QueryRow(query, id).Scan(
		&detail.Bon.ID, &typeBonStr, &detail.Bon.Numero, &detail.Bon.SequenceNumero, &detail.Bon.Statut,
		&detail.Bon.DateBon, &detail.Bon.ReferenceExterne, &detail.Bon.FactureID,
		&detail.Bon.EntrepriseDenomination, &detail.Bon.EntrepriseAdresse, &detail.Bon.EntrepriseTelephone, &detail.Bon.EntrepriseEmail,
		&detail.Bon.ClientID, &detail.Bon.ClientNom, &detail.Bon.ClientAdresse, &detail.Bon.ClientCodePostal,
		&detail.Bon.ClientVille, &detail.Bon.ClientTelephone, &detail.Bon.TotalHT, &detail.Bon.TotalTVA, &detail.Bon.TotalTTC,
		&detail.Bon.Notes, &createdStr, &updatedStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return detail, fmt.Errorf("bon introuvable: %s", id)
		}
		return detail, fmt.Errorf("query bon: %w", err)
	}

	detail.Bon.TypeBon = TypeBon(typeBonStr)
	detail.Bon.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	detail.Bon.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)

	lignesQuery := `SELECT id, bon_id, produit_id, code_barres, reference, designation,
		quantite, prix_unitaire_ht, taux_tva, montant_ht, montant_tva, montant_ttc, ordre,
		created_at, updated_at
		FROM lignes_bon WHERE bon_id = ? ORDER BY ordre ASC`

	rows, err := s.db.Query(lignesQuery, id)
	if err != nil {
		return detail, fmt.Errorf("query lignes bon: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l LigneBon
		var lCreatedStr, lUpdatedStr string
		err := rows.Scan(
			&l.ID, &l.BonID, &l.ProduitID, &l.CodeBarres, &l.Reference, &l.Designation,
			&l.Quantite, &l.PrixUnitaireHT, &l.TauxTVA, &l.MontantHT, &l.MontantTVA, &l.MontantTTC, &l.Ordre,
			&lCreatedStr, &lUpdatedStr,
		)
		if err != nil {
			return detail, fmt.Errorf("scan ligne bon: %w", err)
		}
		l.CreatedAt, _ = time.Parse(time.RFC3339, lCreatedStr)
		l.UpdatedAt, _ = time.Parse(time.RFC3339, lUpdatedStr)
		detail.Lignes = append(detail.Lignes, l)
	}

	return detail, nil
}

func (s *Store) ListBons(filtre FiltreBons) (PageBons, error) {
	page := filtre.Page
	if page < 1 {
		page = 1
	}
	limit := filtre.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	whereClauses := []string{"1=1"}
	var args []any

	if strings.TrimSpace(filtre.Recherche) != "" {
		recherche := "%" + strings.TrimSpace(filtre.Recherche) + "%"
		whereClauses = append(whereClauses, "(numero LIKE ? OR client_nom LIKE ? OR reference_externe LIKE ?)")
		args = append(args, recherche, recherche, recherche)
	}

	if filtre.TypeBon != nil && *filtre.TypeBon != "" {
		whereClauses = append(whereClauses, "type_bon = ?")
		args = append(args, string(*filtre.TypeBon))
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM bons %s", whereSQL)
	if err := s.db.QueryRow(countSQL, args...).Scan(&total); err != nil {
		return PageBons{}, fmt.Errorf("count bons: %w", err)
	}

	querySQL := fmt.Sprintf(`SELECT id, type_bon, numero, sequence_numero, statut, date_bon, reference_externe,
		facture_id, entreprise_denomination, entreprise_adresse, entreprise_telephone, entreprise_email,
		client_id, client_nom, client_adresse, client_code_postal, client_ville, client_telephone,
		total_ht, total_tva, total_ttc, notes, created_at, updated_at
		FROM bons %s ORDER BY sequence_numero DESC LIMIT ? OFFSET ?`, whereSQL)

	queryArgs := append(args, limit, offset)
	rows, err := s.db.Query(querySQL, queryArgs...)
	if err != nil {
		return PageBons{}, fmt.Errorf("query bons: %w", err)
	}
	defer rows.Close()

	var resultats []Bon
	for rows.Next() {
		var b Bon
		var createdStr, updatedStr string
		var typeBonStr string
		err := rows.Scan(
			&b.ID, &typeBonStr, &b.Numero, &b.SequenceNumero, &b.Statut, &b.DateBon, &b.ReferenceExterne,
			&b.FactureID, &b.EntrepriseDenomination, &b.EntrepriseAdresse, &b.EntrepriseTelephone, &b.EntrepriseEmail,
			&b.ClientID, &b.ClientNom, &b.ClientAdresse, &b.ClientCodePostal, &b.ClientVille, &b.ClientTelephone,
			&b.TotalHT, &b.TotalTVA, &b.TotalTTC, &b.Notes, &createdStr, &updatedStr,
		)
		if err != nil {
			return PageBons{}, fmt.Errorf("scan bon: %w", err)
		}
		b.TypeBon = TypeBon(typeBonStr)
		b.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
		b.UpdatedAt, _ = time.Parse(time.RFC3339, updatedStr)
		resultats = append(resultats, b)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PageBons{
		Bons:       resultats,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
