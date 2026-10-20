package parametres

import "time"

type ParametresEntreprise struct {
	ID                        string    `json:"id"`
	Denomination              string    `json:"denomination"`
	FormeJuridique            string    `json:"formeJuridique"`
	Capital                   int64     `json:"capital"`
	Adresse                   string    `json:"adresse"`
	CodePostal                string    `json:"codePostal"`
	Ville                     string    `json:"ville"`
	Siren                     string    `json:"siren"`
	Siret                     string    `json:"siret"`
	RCS                       string    `json:"rcs"`
	NumeroTVA                 string    `json:"numeroTva"`
	Telephone                 string    `json:"telephone"`
	Email                     string    `json:"email"`
	SiteWeb                   string    `json:"siteWeb"`
	LogoPath                  string    `json:"logoPath"`
	Devise                    string    `json:"devise"`
	TauxTVADefaut             int64     `json:"tauxTvaDefaut"`
	TauxTVADisponibles        string    `json:"tauxTvaDisponibles"`
	MentionsLegales           string    `json:"mentionsLegales"`
	FormatNumerotationFacture string    `json:"formatNumerotationFacture"`
	FormatNumerotationBon     string    `json:"formatNumerotationBon"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}

type ModificationParametres struct {
	Denomination              string `json:"denomination"`
	FormeJuridique            string `json:"formeJuridique"`
	Capital                   int64  `json:"capital"`
	Adresse                   string `json:"adresse"`
	CodePostal                string `json:"codePostal"`
	Ville                     string `json:"ville"`
	Siren                     string `json:"siren"`
	Siret                     string `json:"siret"`
	RCS                       string `json:"rcs"`
	NumeroTVA                 string `json:"numeroTva"`
	Telephone                 string `json:"telephone"`
	Email                     string `json:"email"`
	SiteWeb                   string `json:"siteWeb"`
	LogoPath                  string `json:"logoPath"`
	Devise                    string `json:"devise"`
	TauxTVADefaut             int64  `json:"tauxTvaDefaut"`
	TauxTVADisponibles        string `json:"tauxTvaDisponibles"`
	MentionsLegales           string `json:"mentionsLegales"`
	FormatNumerotationFacture string `json:"formatNumerotationFacture"`
	FormatNumerotationBon     string `json:"formatNumerotationBon"`
}
