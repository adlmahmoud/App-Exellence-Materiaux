package drivers

import (
	"context"
	"exellence-materiaux/Payement_carte/models"
)

type TPEDriver interface {
	InitierPaiement(ctx context.Context, req models.DemandePaiement) (models.ReponsePaiement, error)
	TesterConnexion(ctx context.Context) (models.ResultatTestConnexion, error)
	Fermer() error
}
