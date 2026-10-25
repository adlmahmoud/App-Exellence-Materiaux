package drivers

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"

	"exellence-materiaux/Payement_carte/models"
)

type SerieDriver struct {
	config models.ConfigTPE
}

func NewSerieDriver(config models.ConfigTPE) *SerieDriver {
	return &SerieDriver{config: config}
}

func (s *SerieDriver) ouvrirPort() (*os.File, error) {
	port := s.config.PortCOM
	if port == "" {
		port = "COM1"
	}
	path := `\\.\` + port

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	h, err := windows.CreateFile(
		pathPtr,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le port %s: %w", port, err)
	}

	var dcb windows.DCB
	if err := windows.GetCommState(h, &dcb); err == nil {
		baud := uint32(s.config.BaudRate)
		if baud == 0 {
			baud = 9600
		}
		dcb.BaudRate = baud
		dcb.ByteSize = 8
		dcb.Parity = windows.NOPARITY
		dcb.StopBits = windows.ONESTOPBIT
		_ = windows.SetCommState(h, &dcb)
	}

	var timeouts windows.CommTimeouts
	timeouts.ReadIntervalTimeout = 100
	timeouts.ReadTotalTimeoutMultiplier = 10
	timeouts.ReadTotalTimeoutConstant = 1000
	timeouts.WriteTotalTimeoutMultiplier = 10
	timeouts.WriteTotalTimeoutConstant = 1000
	_ = windows.SetCommTimeouts(h, &timeouts)

	return os.NewFile(uintptr(h), path), nil
}

func (s *SerieDriver) InitierPaiement(ctx context.Context, req models.DemandePaiement) (models.ReponsePaiement, error) {
	f, err := s.ouvrirPort()
	if err != nil {
		return models.ReponsePaiement{
			Statut:          models.StatutErreur,
			CodeReponse:     "ER_PORT",
			Message:         err.Error(),
			MontantCentimes: req.MontantCentimes,
			MontantEuros:    req.MontantEuros,
			DateHeure:       time.Now().Format("2006-01-02 15:04:05"),
		}, err
	}
	defer f.Close()

	timeout := time.Duration(s.config.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	return dialoguerConcert(ctx, f, req.NumeroCaisse, req.MontantCentimes, timeout)
}

func (s *SerieDriver) TesterConnexion(ctx context.Context) (models.ResultatTestConnexion, error) {
	f, err := s.ouvrirPort()
	if err != nil {
		return models.ResultatTestConnexion{
			Succes:  false,
			Message: fmt.Sprintf("Échec d'ouverture du port %s : %v", s.config.PortCOM, err),
			Mode:    "serie",
			Details: "Vérifiez que le câble USB du TPE est bien branché et que le pilote est installé.",
		}, err
	}
	defer f.Close()

	return models.ResultatTestConnexion{
		Succes:  true,
		Message: fmt.Sprintf("Port série %s ouvert et accessible avec succès", s.config.PortCOM),
		Mode:    "serie",
		Details: fmt.Sprintf("BaudRate: %d, Prêt pour protocole Concert", s.config.BaudRate),
	}, nil
}

func (s *SerieDriver) Fermer() error {
	return nil
}
