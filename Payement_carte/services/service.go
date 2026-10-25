package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sys/windows"

	"exellence-materiaux/Payement_carte/drivers"
	"exellence-materiaux/Payement_carte/models"
)

type TPEService struct {
	mu           sync.Mutex
	dataDir      string
	config       models.ConfigTPE
	activeCancel context.CancelFunc
	actif        bool
}

func NewTPEService(dataDir string) *TPEService {
	s := &TPEService{
		dataDir: dataDir,
		config: models.ConfigTPE{
			Mode:               "simulateur",
			PortCOM:            "COM1",
			BaudRate:           9600,
			AdresseIP:          "192.168.1.50",
			PortReseau:         8888,
			TimeoutSec:         45,
			SimulateurScenario: "succes",
			SimulateurDelaiSec: 2,
		},
	}
	s.chargerConfig()
	return s
}

func (s *TPEService) cheminConfig() string {
	return filepath.Join(s.dataDir, "config_tpe.json")
}

func (s *TPEService) chargerConfig() {
	p := s.cheminConfig()
	data, err := os.ReadFile(p)
	if err != nil {
		return
	}
	var cfg models.ConfigTPE
	if err := json.Unmarshal(data, &cfg); err == nil {
		if cfg.Mode != "" {
			s.config = cfg
		}
	}
}

func (s *TPEService) SauvegarderConfig(cfg models.ConfigTPE) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cfg.Mode == "" {
		cfg.Mode = "simulateur"
	}
	if cfg.TimeoutSec <= 0 {
		cfg.TimeoutSec = 45
	}
	if cfg.SimulateurDelaiSec <= 0 {
		cfg.SimulateurDelaiSec = 2
	}
	if cfg.SimulateurScenario == "" {
		cfg.SimulateurScenario = "succes"
	}

	s.config = cfg

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("sérialiser configuration TPE: %w", err)
	}

	p := s.cheminConfig()
	if err := os.WriteFile(p, data, 0644); err != nil {
		return fmt.Errorf("enregistrer configuration TPE: %w", err)
	}

	return nil
}

func (s *TPEService) GetConfig() models.ConfigTPE {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config
}

func (s *TPEService) getDriver() drivers.TPEDriver {
	switch s.config.Mode {
	case "serie":
		return drivers.NewSerieDriver(s.config)
	case "reseau":
		return drivers.NewReseauDriver(s.config)
	default:
		return drivers.NewSimulateurDriver(s.config)
	}
}

func (s *TPEService) ListerPortsCOM() []string {
	var disponibles []string
	for i := 1; i <= 32; i++ {
		nomPort := fmt.Sprintf("COM%d", i)
		chemin := `\\.\` + nomPort
		ptr, err := windows.UTF16PtrFromString(chemin)
		if err != nil {
			continue
		}
		h, err := windows.CreateFile(
			ptr,
			windows.GENERIC_READ|windows.GENERIC_WRITE,
			0,
			nil,
			windows.OPEN_EXISTING,
			windows.FILE_ATTRIBUTE_NORMAL,
			0,
		)
		if err == nil {
			disponibles = append(disponibles, nomPort)
			_ = windows.CloseHandle(h)
		} else if err == windows.ERROR_ACCESS_DENIED {
			disponibles = append(disponibles, nomPort)
		}
	}
	return disponibles
}

func (s *TPEService) TesterConnexion() models.ResultatTestConnexion {
	s.mu.Lock()
	cfg := s.config
	s.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var driver drivers.TPEDriver
	switch cfg.Mode {
	case "serie":
		driver = drivers.NewSerieDriver(cfg)
	case "reseau":
		driver = drivers.NewReseauDriver(cfg)
	default:
		driver = drivers.NewSimulateurDriver(cfg)
	}

	res, err := driver.TesterConnexion(ctx)
	if err != nil && res.Message == "" {
		res.Succes = false
		res.Message = err.Error()
		res.Mode = cfg.Mode
	}
	return res
}

func (s *TPEService) LancerPaiement(montantEuros float64, reference string) models.ReponsePaiement {
	s.mu.Lock()
	if s.actif {
		s.mu.Unlock()
		return models.ReponsePaiement{
			Statut:          models.StatutErreur,
			CodeReponse:     "BUSY",
			Message:         "Une transaction TPE est déjà en cours",
			MontantEuros:    montantEuros,
			MontantCentimes: int64(math.Round(montantEuros * 100)),
			DateHeure:       time.Now().Format("2006-01-02 15:04:05"),
		}
	}

	cfg := s.config
	timeoutSec := cfg.TimeoutSec
	if timeoutSec <= 0 {
		timeoutSec = 60
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	s.activeCancel = cancel
	s.actif = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.actif = false
		s.activeCancel = nil
		s.mu.Unlock()
	}()

	montantCentimes := int64(math.Round(montantEuros * 100))
	req := models.DemandePaiement{
		MontantCentimes: montantCentimes,
		MontantEuros:    montantEuros,
		Reference:       reference,
		NumeroCaisse:    1,
	}

	driver := s.getDriver()
	resp, err := driver.InitierPaiement(ctx, req)
	if err != nil && resp.Statut == "" {
		resp.Statut = models.StatutErreur
		resp.Message = err.Error()
		resp.MontantEuros = montantEuros
		resp.MontantCentimes = montantCentimes
		resp.DateHeure = time.Now().Format("2006-01-02 15:04:05")
	}

	return resp
}

func (s *TPEService) AnnulerPaiementEnCours() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeCancel != nil {
		s.activeCancel()
	}
}

func (s *TPEService) SetSimulateurScenario(scenario string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.SimulateurScenario = scenario
}
