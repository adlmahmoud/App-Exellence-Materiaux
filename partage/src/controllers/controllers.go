package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"exellence-materiaux/internal/auth"
	"exellence-materiaux/internal/factures"
	"exellence-materiaux/internal/parametres"
	"exellence-materiaux/internal/produits"
	"exellence-materiaux/internal/stock"
	"exellence-materiaux/partage/src/helpers"
	"exellence-materiaux/partage/src/models"
	"exellence-materiaux/partage/src/services"
)

type ControllerManager struct {
	services *services.ServiceManager
}

func NewControllerManager(sm *services.ServiceManager) *ControllerManager {
	return &ControllerManager{services: sm}
}

func (c *ControllerManager) CheckAuth(r *http.Request, allowedRoles ...string) (string, error) {
	token := helpers.ExtractBearerToken(r)
	if token == "" {
		return "", fmt.Errorf("token d'authentification manquant")
	}

	role, valid := c.services.Auth.ValidateToken(token)
	if !valid {
		return "", fmt.Errorf("session invalide ou expirée")
	}

	if len(allowedRoles) == 0 {
		return role, nil
	}

	for _, allowed := range allowedRoles {
		if role == allowed {
			return role, nil
		}
	}
	return "", fmt.Errorf("accès refusé : rôle insuffisant (%s)", role)
}

func (c *ControllerManager) AuthStatus(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	configured, err := c.services.Auth.IsConfigured()
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, map[string]bool{"configured": configured}, "statut d'authentification")
}

func (c *ControllerManager) AuthSetup(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	var req models.SetupRequest
	if err := helpers.ReadJSON(r, &req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	if err := c.services.Auth.Setup(req.AdminPassword, req.UserPassword); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusCreated, nil, "mots de passe configurés avec succès")
}

func (c *ControllerManager) AuthLogin(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	var req models.LoginRequest
	if err := helpers.ReadJSON(r, &req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "corps de requête JSON invalide")
		return
	}

	token, expiresAt, err := c.services.Auth.Login(req.Role, req.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, "identifiants invalides")
		return
	}

	helpers.WriteSuccess(w, http.StatusOK, models.LoginResponse{
		Token:     token,
		Role:      req.Role,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, "connexion réussie")
}

func (c *ControllerManager) AuthLogout(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	token := helpers.ExtractBearerToken(r)
	if token != "" {
		c.services.Auth.Logout(token)
	}
	helpers.WriteSuccess(w, http.StatusOK, nil, "déconnexion effectuée")
}

func (c *ControllerManager) AuthChangePasswords(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var req models.ChangePasswordRequest
	if err := helpers.ReadJSON(r, &req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données invalides")
		return
	}

	if err := c.services.Auth.ChangePasswords(req.CurrentAdminPassword, req.AdminPassword, req.UserPassword); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, nil, "mots de passe modifiés avec succès")
}

func (c *ControllerManager) ProduitsList(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	recherche := helpers.GetQueryParam(r, "recherche", "")
	page := helpers.GetQueryInt(r, "page", 1)
	limit := helpers.GetQueryInt(r, "limit", 50)
	alerteStock := helpers.GetQueryBool(r, "alerteStock", false)

	res, err := c.services.Produits.List(produits.FiltreProduits{
		Recherche:   recherche,
		Page:        page,
		Limit:       limit,
		AlerteStock: alerteStock,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, res, "")
}

func (c *ControllerManager) ProduitsGet(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	barcode := helpers.GetQueryParam(r, "code", "")

	if barcode != "" {
		p, err := c.services.Produits.GetByCodeBarres(barcode)
		if err != nil {
			helpers.WriteError(w, http.StatusNotFound, "produit introuvable")
			return
		}
		helpers.WriteSuccess(w, http.StatusOK, p, "")
		return
	}

	if id != "" {
		p, err := c.services.Produits.GetByID(id)
		if err != nil {
			helpers.WriteError(w, http.StatusNotFound, "produit introuvable")
			return
		}
		helpers.WriteSuccess(w, http.StatusOK, p, "")
		return
	}

	helpers.WriteError(w, http.StatusBadRequest, "paramètre id ou code requis")
}

func (c *ControllerManager) ProduitsCreate(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var np produits.NouveauProduit
	if err := helpers.ReadJSON(r, &np); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données invalides")
		return
	}

	cree, err := c.services.Produits.Create(np)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusCreated, cree, "produit créé avec succès")
}

func (c *ControllerManager) ProduitsUpdate(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var mp produits.ModificationProduit
	if err := helpers.ReadJSON(r, &mp); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données invalides")
		return
	}

	maj, err := c.services.Produits.Update(mp)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, maj, "produit mis à jour avec succès")
}

func (c *ControllerManager) ProduitsDelete(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	if id == "" {
		var req struct {
			ID string `json:"id"`
		}
		_ = helpers.ReadJSON(r, &req)
		id = req.ID
	}

	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "identifiant manquant")
		return
	}

	if err := c.services.Produits.Delete(id); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, nil, "produit supprimé avec succès")
}

func (c *ControllerManager) FacturesList(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	recherche := helpers.GetQueryParam(r, "recherche", "")
	page := helpers.GetQueryInt(r, "page", 1)
	limit := helpers.GetQueryInt(r, "limit", 50)

	res, err := c.services.Factures.List(factures.FiltreFactures{
		Recherche: recherche,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, res, "")
}

func (c *ControllerManager) FacturesGet(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "identifiant facture manquant")
		return
	}

	detail, err := c.services.Factures.GetByID(id)
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, "facture introuvable")
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, detail, "")
}

func (c *ControllerManager) FacturesCreate(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var nf factures.NouvelleFacture
	if err := helpers.ReadJSON(r, &nf); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données facture invalides")
		return
	}

	detail, err := c.services.Factures.Create(nf)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusCreated, detail, "facture générée avec succès")
}

func (c *ControllerManager) FacturesPDF(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "identifiant facture manquant")
		return
	}

	filePath, filename, err := c.services.Factures.GeneratePDF(id)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer os.Remove(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "erreur lecture pdf")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	_, _ = io.Copy(w, file)
}

func (c *ControllerManager) BonsList(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	recherche := helpers.GetQueryParam(r, "recherche", "")
	page := helpers.GetQueryInt(r, "page", 1)
	limit := helpers.GetQueryInt(r, "limit", 50)

	res, err := c.services.Bons.List(factures.FiltreBons{
		Recherche: recherche,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, res, "")
}

func (c *ControllerManager) BonsGet(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "identifiant bon manquant")
		return
	}

	detail, err := c.services.Bons.GetByID(id)
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, "bon introuvable")
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, detail, "")
}

func (c *ControllerManager) BonsCreate(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var nb factures.NouveauBon
	if err := helpers.ReadJSON(r, &nb); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données bon invalides")
		return
	}

	detail, err := c.services.Bons.Create(nb)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusCreated, detail, "bon créé avec succès")
}

func (c *ControllerManager) BonsPDF(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id := helpers.GetQueryParam(r, "id", "")
	if id == "" {
		helpers.WriteError(w, http.StatusBadRequest, "identifiant bon manquant")
		return
	}

	filePath, filename, err := c.services.Bons.GeneratePDF(id)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer os.Remove(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "erreur lecture pdf")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	_, _ = io.Copy(w, file)
}

func (c *ControllerManager) StockMouvements(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
			helpers.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		page := helpers.GetQueryInt(r, "page", 1)
		limit := helpers.GetQueryInt(r, "limit", 50)
		produitID := helpers.GetQueryParam(r, "produitId", "")
		typeMvt := helpers.GetQueryParam(r, "typeMouvement", "")

		var tmPtr *stock.TypeMouvement
		if typeMvt != "" {
			tm := stock.TypeMouvement(typeMvt)
			tmPtr = &tm
		}

		res, err := c.services.Stock.ListMouvements(stock.FiltreMouvements{
			Page:          page,
			Limit:         limit,
			ProduitID:     produitID,
			TypeMouvement: tmPtr,
		})
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		helpers.WriteSuccess(w, http.StatusOK, res, "")
		return
	}

	if r.Method == http.MethodPost {
		if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
			helpers.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		var sm stock.SaisieMouvement
		if err := helpers.ReadJSON(r, &sm); err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "données mouvement invalides")
			return
		}
		mvt, err := c.services.Stock.CreateMouvement(sm)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		helpers.WriteSuccess(w, http.StatusCreated, mvt, "mouvement enregistré")
		return
	}

	helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
}

func (c *ControllerManager) StockAjuster(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var sa stock.SaisieAjustement
	if err := helpers.ReadJSON(r, &sa); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données ajustement invalides")
		return
	}

	mvt, err := c.services.Stock.Ajuster(sa)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, mvt, "stock ajusté avec succès")
}

func (c *ControllerManager) StockAlertes(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin, auth.RoleUser); err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	alertes, err := c.services.Stock.GetAlertes()
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, alertes, "")
}

func (c *ControllerManager) ParametresGet(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	p, err := c.services.Parametres.Get()
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, p, "")
}

func (c *ControllerManager) ParametresUpdate(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var mp parametres.ModificationParametres
	if err := helpers.ReadJSON(r, &mp); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "données paramètres invalides")
		return
	}

	maj, err := c.services.Parametres.Update(mp)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, maj, "paramètres mis à jour")
}

func (c *ControllerManager) ParametresUploadLogo(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		return
	}

	if _, err := c.CheckAuth(r, auth.RoleAdmin); err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	var req struct {
		ImageBase64 string `json:"imageBase64"`
	}
	if err := helpers.ReadJSON(r, &req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "image base64 requise")
		return
	}

	raw := req.ImageBase64
	if strings.Contains(raw, ",") {
		raw = strings.Split(raw, ",")[1]
	}

	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "décodage base64 échoué")
		return
	}

	path, err := c.services.Parametres.SaveLogo(data)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, map[string]string{"path": path}, "logo mis à jour")
}

func (c *ControllerManager) Health(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, models.ServerStatus{
		Status:    "online",
		Version:   "1.0.0",
		Database:  "exellence.db (SQLite WAL)",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, "serveur actif")
}

func (c *ControllerManager) APropos(w http.ResponseWriter, r *http.Request) {
	if helpers.EnableCORS(w, r) {
		return
	}
	helpers.WriteSuccess(w, http.StatusOK, map[string]string{
		"nom":     "Exellence Matériaux API",
		"version": "1.0.0",
		"annee":   "2026",
		"auteurs": "Adel Mahmoud et Sidou",
	}, "")
}
