package models

type StatutTransaction string

const (
	StatutEnCours StatutTransaction = "EN_COURS"
	StatutValide  StatutTransaction = "VALIDE"
	StatutRefuse  StatutTransaction = "REFUSE"
	StatutAnnule  StatutTransaction = "ANNULE"
	StatutErreur  StatutTransaction = "ERREUR"
)

type ModeConnexion string

const (
	ModeSimulateur ModeConnexion = "simulateur"
	ModeSerie      ModeConnexion = "serie"
	ModeReseau     ModeConnexion = "reseau"
)

type ConfigTPE struct {
	Mode               string `json:"mode"`
	PortCOM            string `json:"portCom"`
	BaudRate           int    `json:"baudRate"`
	AdresseIP          string `json:"adresseIp"`
	PortReseau         int    `json:"portReseau"`
	TimeoutSec         int    `json:"timeoutSec"`
	SimulateurScenario string `json:"simulateurScenario"`
	SimulateurDelaiSec int    `json:"simulateurDelaiSec"`
}

type DemandePaiement struct {
	MontantCentimes int64   `json:"montantCentimes"`
	MontantEuros    float64 `json:"montantEuros"`
	Reference       string  `json:"reference"`
	NumeroCaisse    int     `json:"numeroCaisse"`
}

type ReponsePaiement struct {
	Statut             StatutTransaction `json:"statut"`
	CodeReponse        string            `json:"codeReponse"`
	Message            string            `json:"message"`
	MontantCentimes    int64             `json:"montantCentimes"`
	MontantEuros       float64           `json:"montantEuros"`
	NumeroAutorisation string            `json:"numeroAutorisation"`
	NumeroTransaction  string            `json:"numeroTransaction"`
	TypeCarte          string            `json:"typeCarte"`
	DateHeure          string            `json:"dateHeure"`
}

type ResultatTestConnexion struct {
	Succes  bool   `json:"succes"`
	Message string `json:"message"`
	Mode    string `json:"mode"`
	Details string `json:"details"`
}
