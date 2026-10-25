package drivers

import (
	"context"
	"fmt"
	"net"
	"time"

	"exellence-materiaux/Payement_carte/models"
)

type ReseauDriver struct {
	config models.ConfigTPE
}

func NewReseauDriver(config models.ConfigTPE) *ReseauDriver {
	return &ReseauDriver{config: config}
}

func (r *ReseauDriver) InitierPaiement(ctx context.Context, req models.DemandePaiement) (models.ReponsePaiement, error) {
	host := r.config.AdresseIP
	port := r.config.PortReseau
	if host == "" {
		host = "127.0.0.1"
	}
	if port <= 0 {
		port = 8888
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	timeout := time.Duration(r.config.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return models.ReponsePaiement{
			Statut:          models.StatutErreur,
			CodeReponse:     "ER_CONN",
			Message:         fmt.Sprintf("Impossible de joindre le TPE réseau à l'adresse %s : %v", addr, err),
			MontantCentimes: req.MontantCentimes,
			MontantEuros:    req.MontantEuros,
			DateHeure:       time.Now().Format("2006-01-02 15:04:05"),
		}, err
	}
	defer conn.Close()

	return dialoguerConcert(ctx, conn, req.NumeroCaisse, req.MontantCentimes, timeout)
}

func (r *ReseauDriver) TesterConnexion(ctx context.Context) (models.ResultatTestConnexion, error) {
	host := r.config.AdresseIP
	port := r.config.PortReseau
	if host == "" {
		host = "127.0.0.1"
	}
	if port <= 0 {
		port = 8888
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return models.ResultatTestConnexion{
			Succes:  false,
			Message: fmt.Sprintf("Échec de connexion au TPE réseau (%s) : %v", addr, err),
			Mode:    "reseau",
			Details: "Vérifiez que le TPE est allumé, connecté au même réseau et que le port est correct.",
		}, err
	}
	defer conn.Close()

	return models.ResultatTestConnexion{
		Succes:  true,
		Message: fmt.Sprintf("TPE réseau détecté et connecté avec succès à %s", addr),
		Mode:    "reseau",
		Details: "Communication TCP active",
	}, nil
}

func (r *ReseauDriver) Fermer() error {
	return nil
}
