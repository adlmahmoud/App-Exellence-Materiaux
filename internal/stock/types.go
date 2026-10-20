package stock

import "time"

type TypeMouvement string

const (
	MouvementEntree     TypeMouvement = "entree"
	MouvementSortie     TypeMouvement = "sortie"
	MouvementAjustement TypeMouvement = "ajustement"
	MouvementVente      TypeMouvement = "vente"
	MouvementRetour     TypeMouvement = "retour"
	MouvementReception  TypeMouvement = "reception"
)

type MouvementStock struct {
	ID                string        `json:"id"`
	ProduitID         string        `json:"produitId"`
	TypeMouvement     TypeMouvement `json:"typeMouvement"`
	Quantite          int64         `json:"quantite"`
	StockAvant        int64         `json:"stockAvant"`
	StockApres        int64         `json:"stockApres"`
	PrixUnitaireHT    int64         `json:"prixUnitaireHt"`
	ReferenceDocument string        `json:"referenceDocument"`
	Motif             string        `json:"motif"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

type MouvementDetail struct {
	MouvementStock
	CodeBarres  string `json:"codeBarres"`
	Reference   string `json:"reference"`
	Designation string `json:"designation"`
}

type SaisieMouvement struct {
	ProduitID         string        `json:"produitId"`
	TypeMouvement     TypeMouvement `json:"typeMouvement"`
	Quantite          int64         `json:"quantite"`
	PrixUnitaireHT    int64         `json:"prixUnitaireHt"`
	ReferenceDocument string        `json:"referenceDocument"`
	Motif             string        `json:"motif"`
}

type SaisieAjustement struct {
	ProduitID    string `json:"produitId"`
	NouveauStock int64  `json:"nouveauStock"`
	Motif        string `json:"motif"`
}

type FiltreMouvements struct {
	ProduitID     string         `json:"produitId"`
	TypeMouvement *TypeMouvement `json:"typeMouvement"`
	DateDebut     string         `json:"dateDebut"`
	DateFin       string         `json:"dateFin"`
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
}

type PageMouvements struct {
	Mouvements []MouvementDetail `json:"mouvements"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

type AlerteStock struct {
	ProduitID    string `json:"produitId"`
	CodeBarres   string `json:"codeBarres"`
	Reference    string `json:"reference"`
	Designation  string `json:"designation"`
	StockActuel  int64  `json:"stockActuel"`
	StockMinimum int64  `json:"stockMinimum"`
	Unite        string `json:"unite"`
}
