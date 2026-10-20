package factures

import "time"

type StatutFacture string

const (
	StatutBrouillon StatutFacture = "brouillon"
	StatutValidee   StatutFacture = "validee"
	StatutPayee     StatutFacture = "payee"
	StatutAnnulee   StatutFacture = "annulee"
)

type ModeReglement string

const (
	ReglementEspeces  ModeReglement = "especes"
	ReglementCB       ModeReglement = "cb"
	ReglementVirement ModeReglement = "virement"
	ReglementCheque   ModeReglement = "cheque"
)

type TypeBon string

const (
	BonLivraison TypeBon = "livraison"
	BonCommande  TypeBon = "commande"
	BonDevis     TypeBon = "devis"
	BonReception TypeBon = "reception"
)

type Client struct {
	ID         string    `json:"id"`
	Nom        string    `json:"nom"`
	Adresse    string    `json:"adresse"`
	CodePostal string    `json:"codePostal"`
	Ville      string    `json:"ville"`
	Telephone  string    `json:"telephone"`
	Email      string    `json:"email"`
	SirenSiret string    `json:"sirenSiret"`
	NumeroTVA  string    `json:"numeroTva"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type NouveauClient struct {
	Nom        string `json:"nom"`
	Adresse    string `json:"adresse"`
	CodePostal string `json:"codePostal"`
	Ville      string `json:"ville"`
	Telephone  string `json:"telephone"`
	Email      string `json:"email"`
	SirenSiret string `json:"sirenSiret"`
	NumeroTVA  string `json:"numeroTva"`
	Notes      string `json:"notes"`
}

type Facture struct {
	ID                       string        `json:"id"`
	Numero                   string        `json:"numero"`
	SequenceNumero           int64         `json:"sequenceNumero"`
	Statut                   StatutFacture `json:"statut"`
	ModeReglement            ModeReglement `json:"modeReglement"`
	DateFacture              string        `json:"dateFacture"`
	DateEcheance             string        `json:"dateEcheance"`
	ConditionsReglement      string        `json:"conditionsReglement"`
	EntrepriseDenomination   string        `json:"entrepriseDenomination"`
	EntrepriseFormeJuridique string        `json:"entrepriseFormeJuridique"`
	EntrepriseCapital        int64         `json:"entrepriseCapital"`
	EntrepriseAdresse        string        `json:"entrepriseAdresse"`
	EntrepriseCodePostal     string        `json:"entrepriseCodePostal"`
	EntrepriseVille          string        `json:"entrepriseVille"`
	EntrepriseSiren          string        `json:"entrepriseSiren"`
	EntrepriseSiret          string        `json:"entrepriseSiret"`
	EntrepriseRCS            string        `json:"entrepriseRcs"`
	EntrepriseNumeroTVA      string        `json:"entrepriseNumeroTva"`
	EntrepriseTelephone      string        `json:"entrepriseTelephone"`
	EntrepriseEmail          string        `json:"entrepriseEmail"`
	MentionsLegales          string        `json:"mentionsLegales"`
	Devise                   string        `json:"devise"`
	ClientID                 *string       `json:"clientId"`
	ClientNom                string        `json:"clientNom"`
	ClientAdresse            string        `json:"clientAdresse"`
	ClientCodePostal         string        `json:"clientCodePostal"`
	ClientVille              string        `json:"clientVille"`
	ClientTelephone          string        `json:"clientTelephone"`
	ClientEmail              string        `json:"clientEmail"`
	ClientSirenSiret         string        `json:"clientSirenSiret"`
	ClientNumeroTVA          string        `json:"clientNumeroTva"`
	TotalHT                  int64         `json:"totalHt"`
	TotalTVA                 int64         `json:"totalTva"`
	TotalTTC                 int64         `json:"totalTtc"`
	MontantRegle             int64         `json:"montantRegle"`
	ResteAPayer              int64         `json:"resteAPayer"`
	Notes                    string        `json:"notes"`
	CreatedAt                time.Time     `json:"createdAt"`
	UpdatedAt                time.Time     `json:"updatedAt"`
}

type LigneFacture struct {
	ID             string    `json:"id"`
	FactureID      string    `json:"factureId"`
	ProduitID      *string   `json:"produitId"`
	CodeBarres     string    `json:"codeBarres"`
	Reference      string    `json:"reference"`
	Designation    string    `json:"designation"`
	Quantite       int64     `json:"quantite"`
	PrixUnitaireHT int64     `json:"prixUnitaireHt"`
	TauxTVA        int64     `json:"tauxTva"`
	MontantHT      int64     `json:"montantHt"`
	MontantTVA     int64     `json:"montantTva"`
	MontantTTC     int64     `json:"montantTtc"`
	Ordre          int       `json:"ordre"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type DetailFacture struct {
	Facture Facture        `json:"facture"`
	Lignes  []LigneFacture `json:"lignes"`
}

type SaisieLigneFacture struct {
	ProduitID      *string `json:"produitId"`
	CodeBarres     string  `json:"codeBarres"`
	Reference      string  `json:"reference"`
	Designation    string  `json:"designation"`
	Quantite       int64   `json:"quantite"`
	PrixUnitaireHT int64   `json:"prixUnitaireHt"`
	TauxTVA        int64   `json:"tauxTva"`
}

type NouvelleFacture struct {
	ClientID            *string              `json:"clientId"`
	ClientNom           string               `json:"clientNom"`
	ClientAdresse       string               `json:"clientAdresse"`
	ClientSirenSiret    string               `json:"clientSirenSiret"`
	ModeReglement       ModeReglement        `json:"modeReglement"`
	DateFacture         string               `json:"dateFacture"`
	DateEcheance        string               `json:"dateEcheance"`
	ConditionsReglement string               `json:"conditionsReglement"`
	MontantRegle        int64                `json:"montantRegle"`
	Notes               string               `json:"notes"`
	Lignes              []SaisieLigneFacture `json:"lignes"`
}

type FiltreFactures struct {
	Recherche string         `json:"recherche"`
	Statut    *StatutFacture `json:"statut"`
	ClientID  string         `json:"clientId"`
	DateDebut string         `json:"dateDebut"`
	DateFin   string         `json:"dateFin"`
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
}

type PageFactures struct {
	Factures   []Facture `json:"factures"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}

type Bon struct {
	ID                     string    `json:"id"`
	TypeBon                TypeBon   `json:"typeBon"`
	Numero                 string    `json:"numero"`
	SequenceNumero         int64     `json:"sequenceNumero"`
	Statut                 string    `json:"statut"`
	DateBon                string    `json:"dateBon"`
	ReferenceExterne       string    `json:"referenceExterne"`
	FactureID              *string   `json:"factureId"`
	EntrepriseDenomination string    `json:"entrepriseDenomination"`
	EntrepriseAdresse      string    `json:"entrepriseAdresse"`
	EntrepriseTelephone    string    `json:"entrepriseTelephone"`
	EntrepriseEmail        string    `json:"entrepriseEmail"`
	ClientID               *string   `json:"clientId"`
	ClientNom              string    `json:"clientNom"`
	ClientAdresse          string    `json:"clientAdresse"`
	ClientCodePostal       string    `json:"clientCodePostal"`
	ClientVille            string    `json:"clientVille"`
	ClientTelephone        string    `json:"clientTelephone"`
	TotalHT                int64     `json:"totalHt"`
	TotalTVA               int64     `json:"totalTva"`
	TotalTTC               int64     `json:"totalTtc"`
	Notes                  string    `json:"notes"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type LigneBon struct {
	ID             string    `json:"id"`
	BonID          string    `json:"bonId"`
	ProduitID      *string   `json:"produitId"`
	CodeBarres     string    `json:"codeBarres"`
	Reference      string    `json:"reference"`
	Designation    string    `json:"designation"`
	Quantite       int64     `json:"quantite"`
	PrixUnitaireHT int64     `json:"prixUnitaireHt"`
	TauxTVA        int64     `json:"tauxTva"`
	MontantHT      int64     `json:"montantHt"`
	MontantTVA     int64     `json:"montantTva"`
	MontantTTC     int64     `json:"montantTtc"`
	Ordre          int       `json:"ordre"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type DetailBon struct {
	Bon    Bon        `json:"bon"`
	Lignes []LigneBon `json:"lignes"`
}

type SaisieLigneBon struct {
	ProduitID      *string `json:"produitId"`
	CodeBarres     string  `json:"codeBarres"`
	Reference      string  `json:"reference"`
	Designation    string  `json:"designation"`
	Quantite       int64   `json:"quantite"`
	PrixUnitaireHT int64   `json:"prixUnitaireHt"`
	TauxTVA        int64   `json:"tauxTva"`
}

type NouveauBon struct {
	TypeBon          TypeBon          `json:"typeBon"`
	DateBon          string           `json:"dateBon"`
	ReferenceExterne string           `json:"referenceExterne"`
	ClientID         *string          `json:"clientId"`
	ClientNom        string           `json:"clientNom"`
	Notes            string           `json:"notes"`
	Lignes           []SaisieLigneBon `json:"lignes"`
}

type FiltreBons struct {
	TypeBon   *TypeBon `json:"typeBon"`
	Recherche string   `json:"recherche"`
	Statut    string   `json:"statut"`
	ClientID  string   `json:"clientId"`
	DateDebut string   `json:"dateDebut"`
	DateFin   string   `json:"dateFin"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}

type PageBons struct {
	Bons       []Bon `json:"bons"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}
