package controllers

import (
	"exellence-materiaux/Payement_carte/models"
	"exellence-materiaux/Payement_carte/services"
)

type TPEController struct {
	service *services.TPEService
}

func NewTPEController(service *services.TPEService) *TPEController {
	return &TPEController{service: service}
}

func (c *TPEController) LancerPaiement(montantEuros float64, reference string) models.ReponsePaiement {
	return c.service.LancerPaiement(montantEuros, reference)
}

func (c *TPEController) AnnulerPaiement() {
	c.service.AnnulerPaiementEnCours()
}

func (c *TPEController) TesterConnexion() models.ResultatTestConnexion {
	return c.service.TesterConnexion()
}

func (c *TPEController) GetConfig() models.ConfigTPE {
	return c.service.GetConfig()
}

func (c *TPEController) SauvegarderConfig(cfg models.ConfigTPE) error {
	return c.service.SauvegarderConfig(cfg)
}

func (c *TPEController) ListerPortsCOM() []string {
	return c.service.ListerPortsCOM()
}

func (c *TPEController) ChangerScenarioSimulateur(scenario string) {
	c.service.SetSimulateurScenario(scenario)
}
