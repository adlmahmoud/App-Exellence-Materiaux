package drivers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"exellence-materiaux/Payement_carte/models"
)

type SimulateurDriver struct {
	config models.ConfigTPE
}

func NewSimulateurDriver(config models.ConfigTPE) *SimulateurDriver {
	return &SimulateurDriver{config: config}
}

func (s *SimulateurDriver) InitierPaiement(ctx context.Context, req models.DemandePaiement) (models.ReponsePaiement, error) {
	delai := s.config.SimulateurDelaiSec
	if delai <= 0 {
		delai = 2
	}

	timer := time.NewTimer(time.Duration(delai) * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return models.ReponsePaiement{
			Statut:          models.StatutAnnule,
			CodeReponse:     "99",
			Message:         "Paiement annulé par le logiciel ou délai dépassé",
			MontantCentimes: req.MontantCentimes,
			MontantEuros:    req.MontantEuros,
			DateHeure:       time.Now().Format("2006-01-02 15:04:05"),
		}, ctx.Err()
	case <-timer.C:
	}

	dateStr := time.Now().Format("2006-01-02 15:04:05")
	scenario := s.config.SimulateurScenario
	if scenario == "" {
		scenario = "succes"
	}

	switch scenario {
	case "refus":
		return models.ReponsePaiement{
			Statut:          models.StatutRefuse,
			CodeReponse:     "05",
			Message:         "Paiement refusé par le centre bancaire (Solde insuffisant / Plafond atteint)",
			MontantCentimes: req.MontantCentimes,
			MontantEuros:    req.MontantEuros,
			DateHeure:       dateStr,
		}, nil

	case "annulation":
		return models.ReponsePaiement{
			Statut:          models.StatutAnnule,
			CodeReponse:     "99",
			Message:         "Transaction annulée sur le terminal (Touche rouge / Abandon client)",
			MontantCentimes: req.MontantCentimes,
			MontantEuros:    req.MontantEuros,
			DateHeure:       dateStr,
		}, nil

	case "aleatoire":
		r := rand.Intn(100)
		if r < 15 {
			return models.ReponsePaiement{
				Statut:          models.StatutRefuse,
				CodeReponse:     "05",
				Message:         "Paiement refusé par le centre bancaire (Simulation aléatoire)",
				MontantCentimes: req.MontantCentimes,
				MontantEuros:    req.MontantEuros,
				DateHeure:       dateStr,
			}, nil
		} else if r < 25 {
			return models.ReponsePaiement{
				Statut:          models.StatutAnnule,
				CodeReponse:     "99",
				Message:         "Transaction annulée sur le TPE (Simulation aléatoire)",
				MontantCentimes: req.MontantCentimes,
				MontantEuros:    req.MontantEuros,
				DateHeure:       dateStr,
			}, nil
		}
	}

	numAuth := fmt.Sprintf("AUT-%06d", rand.Intn(999999))
	numTrans := fmt.Sprintf("TRX-%08d", rand.Intn(99999999))

	return models.ReponsePaiement{
		Statut:             models.StatutValide,
		CodeReponse:        "00",
		Message:            "Paiement CB accepté par le terminal",
		MontantCentimes:    req.MontantCentimes,
		MontantEuros:       req.MontantEuros,
		NumeroAutorisation: numAuth,
		NumeroTransaction:  numTrans,
		TypeCarte:          "CB EMV / Sans Contact",
		DateHeure:          dateStr,
	}, nil
}

func (s *SimulateurDriver) TesterConnexion(ctx context.Context) (models.ResultatTestConnexion, error) {
	return models.ResultatTestConnexion{
		Succes:  true,
		Message: "Simulateur TPE opérationnel et prêt pour les tests",
		Mode:    "simulateur",
		Details: fmt.Sprintf("Scénario configuré : %s, délai : %ds", s.config.SimulateurScenario, s.config.SimulateurDelaiSec),
	}, nil
}

func (s *SimulateurDriver) Fermer() error {
	return nil
}
