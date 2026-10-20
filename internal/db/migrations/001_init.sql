CREATE TABLE IF NOT EXISTS parametres_entreprise (
    id TEXT PRIMARY KEY,
    denomination TEXT NOT NULL DEFAULT 'Exellence Materiaux',
    forme_juridique TEXT NOT NULL DEFAULT '',
    capital INTEGER NOT NULL DEFAULT 0,
    adresse TEXT NOT NULL DEFAULT '',
    code_postal TEXT NOT NULL DEFAULT '',
    ville TEXT NOT NULL DEFAULT '',
    siren TEXT NOT NULL DEFAULT '',
    siret TEXT NOT NULL DEFAULT '',
    rcs TEXT NOT NULL DEFAULT '',
    numero_tva TEXT NOT NULL DEFAULT '',
    telephone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    site_web TEXT NOT NULL DEFAULT '',
    logo_path TEXT NOT NULL DEFAULT '',
    devise TEXT NOT NULL DEFAULT 'EUR',
    taux_tva_defaut INTEGER NOT NULL DEFAULT 2000,
    taux_tva_disponibles TEXT NOT NULL DEFAULT '[2000,1000,550,210,0]',
    mentions_legales TEXT NOT NULL DEFAULT '',
    format_numerotation_facture TEXT NOT NULL DEFAULT 'FAC-{YYYY}-{SEQ}',
    format_numerotation_bon TEXT NOT NULL DEFAULT 'BON-{TYPE}-{YYYY}-{SEQ}',
    updated_at TEXT NOT NULL
);

INSERT INTO parametres_entreprise (
    id, denomination, forme_juridique, capital, adresse, code_postal, ville,
    siren, siret, rcs, numero_tva, telephone, email, site_web, logo_path,
    devise, taux_tva_defaut, taux_tva_disponibles, mentions_legales,
    format_numerotation_facture, format_numerotation_bon, updated_at
) VALUES (
    '01DEFAULTCONFIG00000000000',
    'Exellence Materiaux',
    '',
    0,
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    'EUR',
    2000,
    '[2000,1000,550,210,0]',
    '',
    'FAC-{YYYY}-{SEQ}',
    'BON-{TYPE}-{YYYY}-{SEQ}',
    datetime('now')
);

CREATE TABLE IF NOT EXISTS produits (
    id TEXT PRIMARY KEY,
    code_barres TEXT NOT NULL UNIQUE,
    reference TEXT NOT NULL DEFAULT '',
    designation TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    prix_achat_ht INTEGER NOT NULL DEFAULT 0,
    prix_vente_ht INTEGER NOT NULL DEFAULT 0,
    prix_vente_ttc INTEGER NOT NULL DEFAULT 0,
    taux_tva INTEGER NOT NULL DEFAULT 2000,
    stock_actuel INTEGER NOT NULL DEFAULT 0,
    stock_minimum INTEGER NOT NULL DEFAULT 0,
    unite TEXT NOT NULL DEFAULT 'piece',
    actif INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_produits_code_barres ON produits(code_barres);
CREATE INDEX IF NOT EXISTS idx_produits_reference ON produits(reference);
CREATE INDEX IF NOT EXISTS idx_produits_designation ON produits(designation);
CREATE INDEX IF NOT EXISTS idx_produits_stock_alerte ON produits(stock_actuel, stock_minimum) WHERE actif = 1;

CREATE VIRTUAL TABLE IF NOT EXISTS produits_fts USING fts5(
    id UNINDEXED,
    code_barres,
    reference,
    designation,
    description,
    content='produits',
    content_rowid='rowid'
);

CREATE TRIGGER IF NOT EXISTS trg_produits_ai AFTER INSERT ON produits BEGIN
    INSERT INTO produits_fts(rowid, id, code_barres, reference, designation, description)
    VALUES (new.rowid, new.id, new.code_barres, new.reference, new.designation, new.description);
END;

CREATE TRIGGER IF NOT EXISTS trg_produits_ad AFTER DELETE ON produits BEGIN
    INSERT INTO produits_fts(produits_fts, rowid, id, code_barres, reference, designation, description)
    VALUES ('delete', old.rowid, old.id, old.code_barres, old.reference, old.designation, old.description);
END;

CREATE TRIGGER IF NOT EXISTS trg_produits_au AFTER UPDATE ON produits BEGIN
    INSERT INTO produits_fts(produits_fts, rowid, id, code_barres, reference, designation, description)
    VALUES ('delete', old.rowid, old.id, old.code_barres, old.reference, old.designation, old.description);
    INSERT INTO produits_fts(rowid, id, code_barres, reference, designation, description)
    VALUES (new.rowid, new.id, new.code_barres, new.reference, new.designation, new.description);
END;

CREATE TABLE IF NOT EXISTS mouvements_stock (
    id TEXT PRIMARY KEY,
    produit_id TEXT NOT NULL REFERENCES produits(id) ON DELETE RESTRICT,
    type_mouvement TEXT NOT NULL,
    quantite INTEGER NOT NULL,
    stock_avant INTEGER NOT NULL,
    stock_apres INTEGER NOT NULL,
    prix_unitaire_ht INTEGER NOT NULL DEFAULT 0,
    reference_document TEXT NOT NULL DEFAULT '',
    motif TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mouvements_stock_produit ON mouvements_stock(produit_id);
CREATE INDEX IF NOT EXISTS idx_mouvements_stock_created_at ON mouvements_stock(created_at);
CREATE INDEX IF NOT EXISTS idx_mouvements_stock_ref_doc ON mouvements_stock(reference_document);

CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    nom TEXT NOT NULL,
    adresse TEXT NOT NULL DEFAULT '',
    code_postal TEXT NOT NULL DEFAULT '',
    ville TEXT NOT NULL DEFAULT '',
    telephone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    siren_siret TEXT NOT NULL DEFAULT '',
    numero_tva TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_clients_nom ON clients(nom);

CREATE TABLE IF NOT EXISTS sequences (
    nom TEXT PRIMARY KEY,
    derniere_valeur INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS factures (
    id TEXT PRIMARY KEY,
    numero TEXT NOT NULL UNIQUE,
    sequence_numero INTEGER NOT NULL,
    statut TEXT NOT NULL DEFAULT 'validee',
    mode_reglement TEXT NOT NULL DEFAULT 'especes',
    date_facture TEXT NOT NULL,
    date_echeance TEXT NOT NULL DEFAULT '',
    entreprise_denomination TEXT NOT NULL,
    entreprise_forme_juridique TEXT NOT NULL DEFAULT '',
    entreprise_capital INTEGER NOT NULL DEFAULT 0,
    entreprise_adresse TEXT NOT NULL DEFAULT '',
    entreprise_code_postal TEXT NOT NULL DEFAULT '',
    entreprise_ville TEXT NOT NULL DEFAULT '',
    entreprise_siren TEXT NOT NULL DEFAULT '',
    entreprise_siret TEXT NOT NULL DEFAULT '',
    entreprise_rcs TEXT NOT NULL DEFAULT '',
    entreprise_numero_tva TEXT NOT NULL DEFAULT '',
    entreprise_telephone TEXT NOT NULL DEFAULT '',
    entreprise_email TEXT NOT NULL DEFAULT '',
    mentions_legales TEXT NOT NULL DEFAULT '',
    devise TEXT NOT NULL DEFAULT 'EUR',
    client_id TEXT REFERENCES clients(id) ON DELETE SET NULL,
    client_nom TEXT NOT NULL,
    client_adresse TEXT NOT NULL DEFAULT '',
    client_code_postal TEXT NOT NULL DEFAULT '',
    client_ville TEXT NOT NULL DEFAULT '',
    client_telephone TEXT NOT NULL DEFAULT '',
    client_email TEXT NOT NULL DEFAULT '',
    client_siren_siret TEXT NOT NULL DEFAULT '',
    client_numero_tva TEXT NOT NULL DEFAULT '',
    total_ht INTEGER NOT NULL DEFAULT 0,
    total_tva INTEGER NOT NULL DEFAULT 0,
    total_ttc INTEGER NOT NULL DEFAULT 0,
    montant_regle INTEGER NOT NULL DEFAULT 0,
    reste_a_payer INTEGER NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_factures_numero ON factures(numero);
CREATE INDEX IF NOT EXISTS idx_factures_date ON factures(date_facture);
CREATE INDEX IF NOT EXISTS idx_factures_client ON factures(client_id);

CREATE TABLE IF NOT EXISTS lignes_facture (
    id TEXT PRIMARY KEY,
    facture_id TEXT NOT NULL REFERENCES factures(id) ON DELETE CASCADE,
    produit_id TEXT REFERENCES produits(id) ON DELETE SET NULL,
    code_barres TEXT NOT NULL DEFAULT '',
    reference TEXT NOT NULL DEFAULT '',
    designation TEXT NOT NULL,
    quantite INTEGER NOT NULL,
    prix_unitaire_ht INTEGER NOT NULL,
    taux_tva INTEGER NOT NULL,
    montant_ht INTEGER NOT NULL,
    montant_tva INTEGER NOT NULL,
    montant_ttc INTEGER NOT NULL,
    ordre INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lignes_facture_facture ON lignes_facture(facture_id);
CREATE INDEX IF NOT EXISTS idx_lignes_facture_produit ON lignes_facture(produit_id);

CREATE TABLE IF NOT EXISTS bons (
    id TEXT PRIMARY KEY,
    type_bon TEXT NOT NULL,
    numero TEXT NOT NULL UNIQUE,
    sequence_numero INTEGER NOT NULL,
    statut TEXT NOT NULL DEFAULT 'en_cours',
    date_bon TEXT NOT NULL,
    reference_externe TEXT NOT NULL DEFAULT '',
    facture_id TEXT REFERENCES factures(id) ON DELETE SET NULL,
    entreprise_denomination TEXT NOT NULL,
    entreprise_adresse TEXT NOT NULL DEFAULT '',
    entreprise_telephone TEXT NOT NULL DEFAULT '',
    entreprise_email TEXT NOT NULL DEFAULT '',
    client_id TEXT REFERENCES clients(id) ON DELETE SET NULL,
    client_nom TEXT NOT NULL,
    client_adresse TEXT NOT NULL DEFAULT '',
    client_code_postal TEXT NOT NULL DEFAULT '',
    client_ville TEXT NOT NULL DEFAULT '',
    client_telephone TEXT NOT NULL DEFAULT '',
    total_ht INTEGER NOT NULL DEFAULT 0,
    total_tva INTEGER NOT NULL DEFAULT 0,
    total_ttc INTEGER NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bons_numero ON bons(numero);
CREATE INDEX IF NOT EXISTS idx_bons_type ON bons(type_bon);
CREATE INDEX IF NOT EXISTS idx_bons_date ON bons(date_bon);

CREATE TABLE IF NOT EXISTS lignes_bon (
    id TEXT PRIMARY KEY,
    bon_id TEXT NOT NULL REFERENCES bons(id) ON DELETE CASCADE,
    produit_id TEXT REFERENCES produits(id) ON DELETE SET NULL,
    code_barres TEXT NOT NULL DEFAULT '',
    reference TEXT NOT NULL DEFAULT '',
    designation TEXT NOT NULL,
    quantite INTEGER NOT NULL,
    prix_unitaire_ht INTEGER NOT NULL DEFAULT 0,
    taux_tva INTEGER NOT NULL DEFAULT 0,
    montant_ht INTEGER NOT NULL DEFAULT 0,
    montant_tva INTEGER NOT NULL DEFAULT 0,
    montant_ttc INTEGER NOT NULL DEFAULT 0,
    ordre INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_lignes_bon_bon ON lignes_bon(bon_id);
CREATE INDEX IF NOT EXISTS idx_lignes_bon_produit ON lignes_bon(produit_id);
