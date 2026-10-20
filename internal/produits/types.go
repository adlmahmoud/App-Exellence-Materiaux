package produits

import "time"

type Produit struct {
	ID           string    `json:"id"`
	CodeBarres   string    `json:"codeBarres"`
	Reference    string    `json:"reference"`
	Designation  string    `json:"designation"`
	Description  string    `json:"description"`
	PrixAchatHT  int64     `json:"prixAchatHt"`
	PrixVenteHT  int64     `json:"prixVenteHt"`
	PrixVenteTTC int64     `json:"prixVenteTtc"`
	TauxTVA      int64     `json:"tauxTva"`
	StockActuel  int64     `json:"stockActuel"`
	StockMinimum int64     `json:"stockMinimum"`
	Unite        string    `json:"unite"`
	Actif        bool      `json:"actif"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type NouveauProduit struct {
	CodeBarres   string `json:"codeBarres"`
	Reference    string `json:"reference"`
	Designation  string `json:"designation"`
	Description  string `json:"description"`
	PrixAchatHT  int64  `json:"prixAchatHt"`
	PrixVenteHT  int64  `json:"prixVenteHt"`
	PrixVenteTTC int64  `json:"prixVenteTtc"`
	TauxTVA      int64  `json:"tauxTva"`
	StockInitial int64  `json:"stockInitial"`
	StockMinimum int64  `json:"stockMinimum"`
	Unite        string `json:"unite"`
}

type ModificationProduit struct {
	ID           string `json:"id"`
	CodeBarres   string `json:"codeBarres"`
	Reference    string `json:"reference"`
	Designation  string `json:"designation"`
	Description  string `json:"description"`
	PrixAchatHT  int64  `json:"prixAchatHt"`
	PrixVenteHT  int64  `json:"prixVenteHt"`
	PrixVenteTTC int64  `json:"prixVenteTtc"`
	TauxTVA      int64  `json:"tauxTva"`
	StockMinimum int64  `json:"stockMinimum"`
	Unite        string `json:"unite"`
	Actif        bool   `json:"actif"`
}

type FiltreProduits struct {
	Recherche    string `json:"recherche"`
	AlerteStock  bool   `json:"alerteStock"`
	Actif        *bool  `json:"actif"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TriChamp     string `json:"triChamp"`
	TriAscendant bool   `json:"triAscendant"`
}

type PageProduits struct {
	Produits   []Produit `json:"produits"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}
