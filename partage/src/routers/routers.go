package routers

import (
	"net/http"

	"exellence-materiaux/partage/src/controllers"
	"exellence-materiaux/partage/src/helpers"
)

func NewRouter(ctrl *controllers.ControllerManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", ctrl.Health)
	mux.HandleFunc("/api/health", ctrl.Health)
	mux.HandleFunc("/api/apropos", ctrl.APropos)

	mux.HandleFunc("/api/auth/status", ctrl.AuthStatus)
	mux.HandleFunc("/api/auth/setup", ctrl.AuthSetup)
	mux.HandleFunc("/api/auth/login", ctrl.AuthLogin)
	mux.HandleFunc("/api/auth/logout", ctrl.AuthLogout)
	mux.HandleFunc("/api/auth/change-password", ctrl.AuthChangePasswords)

	mux.HandleFunc("/api/produits", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" || r.URL.Query().Get("code") != "" {
				ctrl.ProduitsGet(w, r)
				return
			}
			ctrl.ProduitsList(w, r)
		case http.MethodPost:
			ctrl.ProduitsCreate(w, r)
		case http.MethodPut:
			ctrl.ProduitsUpdate(w, r)
		case http.MethodDelete:
			ctrl.ProduitsDelete(w, r)
		case http.MethodOptions:
			helpers.EnableCORS(w, r)
		default:
			helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		}
	})
	mux.HandleFunc("/api/produits/get", ctrl.ProduitsGet)
	mux.HandleFunc("/api/produits/create", ctrl.ProduitsCreate)
	mux.HandleFunc("/api/produits/update", ctrl.ProduitsUpdate)
	mux.HandleFunc("/api/produits/delete", ctrl.ProduitsDelete)

	mux.HandleFunc("/api/factures", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				ctrl.FacturesGet(w, r)
				return
			}
			ctrl.FacturesList(w, r)
		case http.MethodPost:
			ctrl.FacturesCreate(w, r)
		case http.MethodOptions:
			helpers.EnableCORS(w, r)
		default:
			helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		}
	})
	mux.HandleFunc("/api/factures/get", ctrl.FacturesGet)
	mux.HandleFunc("/api/factures/create", ctrl.FacturesCreate)
	mux.HandleFunc("/api/factures/pdf", ctrl.FacturesPDF)

	mux.HandleFunc("/api/bons", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("id") != "" {
				ctrl.BonsGet(w, r)
				return
			}
			ctrl.BonsList(w, r)
		case http.MethodPost:
			ctrl.BonsCreate(w, r)
		case http.MethodOptions:
			helpers.EnableCORS(w, r)
		default:
			helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		}
	})
	mux.HandleFunc("/api/bons/get", ctrl.BonsGet)
	mux.HandleFunc("/api/bons/create", ctrl.BonsCreate)
	mux.HandleFunc("/api/bons/pdf", ctrl.BonsPDF)

	mux.HandleFunc("/api/stock/mouvements", ctrl.StockMouvements)
	mux.HandleFunc("/api/stock/ajuster", ctrl.StockAjuster)
	mux.HandleFunc("/api/stock/alertes", ctrl.StockAlertes)

	mux.HandleFunc("/api/parametres", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctrl.ParametresGet(w, r)
		case http.MethodPut, http.MethodPost:
			ctrl.ParametresUpdate(w, r)
		case http.MethodOptions:
			helpers.EnableCORS(w, r)
		default:
			helpers.WriteError(w, http.StatusMethodNotAllowed, "méthode non autorisée")
		}
	})
	mux.HandleFunc("/api/parametres/logo", ctrl.ParametresUploadLogo)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.EnableCORS(w, r) {
			return
		}
		mux.ServeHTTP(w, r)
	})
}
