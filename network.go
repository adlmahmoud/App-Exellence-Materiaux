package main

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"exellence-materiaux/internal/factures"
	"exellence-materiaux/internal/parametres"
	"exellence-materiaux/internal/produits"
	"exellence-materiaux/internal/stock"
)

const networkConfigFilename = "network.json"
const defaultNetworkToken = "ExellenceMateriauxNetworkKey2026SecureToken"

type NetworkConfig struct {
	Mode      string `json:"mode"`
	Listen    string `json:"listen,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	Token     string `json:"token,omitempty"`
}

type rpcRequest struct {
	Method       string            `json:"method"`
	Args         []json.RawMessage `json:"args"`
	SessionToken string            `json:"sessionToken,omitempty"`
}

type rpcResponse struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type downloadResult struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

func (a *App) loadNetworkConfig() error {
	configDir := a.dataDir
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		configDir = filepath.Join(userConfigDir, "Gestion de Stock")
	}
	return a.loadNetworkConfigAt(filepath.Join(configDir, networkConfigFilename))
}

func (a *App) loadNetworkConfigAt(configPath string) error {
	a.networkConfig = NetworkConfig{Mode: "standalone"}
	content, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		a.networkConfig = NetworkConfig{Mode: "standalone"}
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", configPath, err)
	}

	var config NetworkConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("parse %s: %w", configPath, err)
	}

	switch config.Mode {
	case "standalone":
	case "server":
		if config.Listen == "" {
			config.Listen = "0.0.0.0:8765"
		}
		if config.Token == "" {
			config.Token = defaultNetworkToken
		} else if len(config.Token) < 32 {
			return fmt.Errorf("%s: server token must be at least 32 characters", configPath)
		}
	case "client":
		parsedURL, err := url.ParseRequestURI(config.ServerURL)
		if err != nil || parsedURL.Scheme != "http" || parsedURL.Host == "" {
			return fmt.Errorf("%s: serverUrl must be a valid http URL on the local network", configPath)
		}
		if config.Token == "" {
			config.Token = defaultNetworkToken
		} else if len(config.Token) < 32 {
			return fmt.Errorf("%s: client token must be at least 32 characters", configPath)
		}
	default:
		return fmt.Errorf("%s: mode must be standalone, server, or client", configPath)
	}
	a.networkConfig = config
	return nil
}

func newNetworkServer(app *App) *http.Server {
	return &http.Server{
		Addr:              app.networkConfig.Listen,
		Handler:           newNetworkHandler(app),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func newNetworkHandler(app *App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "online",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/rpc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		expectedToken := app.networkConfig.Token
		if expectedToken == "" {
			expectedToken = defaultNetworkToken
		}
		if len(token) != len(expectedToken) ||
			subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) != 1 {
			writeRPCError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 16<<20)
		var request rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeRPCError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
			return
		}

		result, err := app.dispatchRPCRequest(request)
		if err != nil {
			writeRPCError(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rpcResponse{Result: result}); err != nil {
			return
		}
	})
	return mux
}

func writeRPCError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(rpcResponse{Error: err.Error()})
}

func decodeRPCArg[T any](request rpcRequest, index int) (T, error) {
	var value T
	if index >= len(request.Args) {
		return value, fmt.Errorf("missing argument %d for %s", index+1, request.Method)
	}
	if err := json.Unmarshal(request.Args[index], &value); err != nil {
		return value, fmt.Errorf("invalid argument %d for %s: %w", index+1, request.Method, err)
	}
	return value, nil
}

func encodeDownload(filename, path string) (downloadResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return downloadResult{}, fmt.Errorf("read generated file: %w", err)
	}
	return downloadResult{
		Filename: filename,
		Data:     base64.StdEncoding.EncodeToString(content),
	}, nil
}

func (a *App) dispatchRPCRequest(request rpcRequest) (any, error) {
	switch request.Method {
	case "GetAuthStatus", "SetupAccessPasswords", "Login", "GetAPropos":
		return a.dispatchRPC(request)
	case "Logout":
		if _, ok := a.rpcSessionRole(request.SessionToken); !ok {
			return nil, fmt.Errorf("session invalide ou expirée")
		}
		a.logoutRPCSession(request.SessionToken)
		return nil, nil
	default:
		role, ok := a.rpcSessionRole(request.SessionToken)
		if !ok {
			return nil, fmt.Errorf("session invalide ou expirée")
		}
		return a.dispatchRPCAs(role, request)
	}
}

func (a *App) dispatchRPC(request rpcRequest) (any, error) {
	switch request.Method {
	case "GetAuthStatus":
		return a.GetAuthStatus()
	case "SetupAccessPasswords":
		adminPassword, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		userPassword, err := decodeRPCArg[string](request, 1)
		if err != nil {
			return nil, err
		}
		return nil, a.SetupAccessPasswords(adminPassword, userPassword)
	case "Login":
		role, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		password, err := decodeRPCArg[string](request, 1)
		if err != nil {
			return nil, err
		}
		return a.createRPCSession(role, password)
	case "Logout":
		a.Logout()
		return nil, nil
	case "GetCurrentRole":
		return a.GetCurrentRole(), nil
	case "ChangeAccessPasswords":
		currentPassword, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		adminPassword, err := decodeRPCArg[string](request, 1)
		if err != nil {
			return nil, err
		}
		userPassword, err := decodeRPCArg[string](request, 2)
		if err != nil {
			return nil, err
		}
		return nil, a.ChangeAccessPasswords(currentPassword, adminPassword, userPassword)
	case "GetAPropos":
		return a.GetAPropos(), nil
	case "GetParametres":
		return a.GetParametres()
	case "UpdateParametres":
		arg, err := decodeRPCArg[parametres.ModificationParametres](request, 0)
		if err != nil {
			return nil, err
		}
		return nil, a.UpdateParametres(arg)
	case "ImporterLogo":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, fmt.Errorf("logo upload requires image data")
		}
		content, err := base64.StdEncoding.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid logo image data: %w", err)
		}
		return a.saveLogo(content)
	case "SauvegarderBase":
		path, err := a.SauvegarderBase("")
		if err != nil {
			return nil, err
		}
		return encodeDownload(filepath.Base(path), path)
	case "CreerProduit":
		arg, err := decodeRPCArg[produits.NouveauProduit](request, 0)
		if err != nil {
			return nil, err
		}
		return a.CreerProduit(arg)
	case "ModifierProduit":
		arg, err := decodeRPCArg[produits.ModificationProduit](request, 0)
		if err != nil {
			return nil, err
		}
		return a.ModifierProduit(arg)
	case "SupprimerProduit":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		return nil, a.SupprimerProduit(arg)
	case "GetProduitParCodeBarres":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		return a.GetProduitParCodeBarres(arg)
	case "ListerProduits":
		arg, err := decodeRPCArg[produits.FiltreProduits](request, 0)
		if err != nil {
			return nil, err
		}
		return a.ListerProduits(arg)
	case "CreerMouvementStock":
		arg, err := decodeRPCArg[stock.SaisieMouvement](request, 0)
		if err != nil {
			return nil, err
		}
		return a.CreerMouvementStock(arg)
	case "AjusterStock":
		arg, err := decodeRPCArg[stock.SaisieAjustement](request, 0)
		if err != nil {
			return nil, err
		}
		return a.AjusterStock(arg)
	case "ListerMouvementsStock":
		arg, err := decodeRPCArg[stock.FiltreMouvements](request, 0)
		if err != nil {
			return nil, err
		}
		return a.ListerMouvementsStock(arg)
	case "SupprimerMouvementStock":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		return nil, a.SupprimerMouvementStock(arg)
	case "ObtenirAlertesStock":
		return a.ObtenirAlertesStock()
	case "CreerFacture":
		arg, err := decodeRPCArg[factures.NouvelleFacture](request, 0)
		if err != nil {
			return nil, err
		}
		return a.CreerFacture(arg)
	case "GetFacture":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		return a.GetFacture(arg)
	case "ListerFactures":
		arg, err := decodeRPCArg[factures.FiltreFactures](request, 0)
		if err != nil {
			return nil, err
		}
		return a.ListerFactures(arg)
	case "ImprimerFacturePDF":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		path, err := a.genererFacturePDF(arg)
		if err != nil {
			return nil, err
		}
		detail, err := a.GetFacture(arg)
		if err != nil {
			return nil, err
		}
		return encodeDownload("Facture_"+detail.Facture.Numero+".pdf", path)
	case "CreerBon":
		arg, err := decodeRPCArg[factures.NouveauBon](request, 0)
		if err != nil {
			return nil, err
		}
		return a.CreerBon(arg)
	case "GetBon":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		return a.GetBon(arg)
	case "ListerBons":
		arg, err := decodeRPCArg[factures.FiltreBons](request, 0)
		if err != nil {
			return nil, err
		}
		return a.ListerBons(arg)
	case "ImprimerBonPDF":
		arg, err := decodeRPCArg[string](request, 0)
		if err != nil {
			return nil, err
		}
		path, err := a.genererBonPDF(arg)
		if err != nil {
			return nil, err
		}
		detail, err := a.GetBon(arg)
		if err != nil {
			return nil, err
		}
		return encodeDownload("Bon_"+detail.Bon.Numero+".pdf", path)
	default:
		return nil, fmt.Errorf("unknown method %q", request.Method)
	}
}

type ServerInfo struct {
	Running   bool     `json:"running"`
	Port      int      `json:"port"`
	IPs       []string `json:"ips"`
	Mode      string   `json:"mode"`
	ServerURL string   `json:"serverUrl"`
}

func (a *App) startPartageServer() {
	a.serverMu.Lock()
	defer a.serverMu.Unlock()

	if a.networkServer != nil {
		return
	}

	if a.database == nil || a.networkConfig.Mode == "client" {
		return
	}

	listenAddr := "0.0.0.0:8765"
	if a.networkConfig.Listen != "" {
		listenAddr = a.networkConfig.Listen
	}

	if a.networkConfig.Token == "" {
		a.networkConfig.Token = defaultNetworkToken
	}

	a.networkServer = &http.Server{
		Addr:              listenAddr,
		Handler:           newNetworkHandler(a),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := a.networkServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			println(fmt.Sprintf("serveur partage error: %v", err))
		}
	}()
}

func (a *App) stopPartageServer() {
	a.serverMu.Lock()
	defer a.serverMu.Unlock()

	if a.networkServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = a.networkServer.Shutdown(ctx)
		a.networkServer = nil
	}
}

func (a *App) GetServerInfo() ServerInfo {
	a.serverMu.Lock()
	running := a.networkServer != nil
	a.serverMu.Unlock()

	ips := getLocalIPs()
	return ServerInfo{
		Running:   running,
		Port:      8765,
		IPs:       ips,
		Mode:      a.networkConfig.Mode,
		ServerURL: a.networkConfig.ServerURL,
	}
}

func (a *App) SetNetworkMode(mode, serverIP string) error {
	configDir := a.dataDir
	if userConfigDir, err := os.UserConfigDir(); err == nil {
		configDir = filepath.Join(userConfigDir, "Gestion de Stock")
	}
	_ = os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, networkConfigFilename)

	if mode == "client" {
		serverIP = strings.TrimSpace(serverIP)
		if serverIP == "" {
			return fmt.Errorf("adresse IP du PC 1 requise")
		}
		if !strings.HasPrefix(serverIP, "http://") && !strings.HasPrefix(serverIP, "https://") {
			serverIP = "http://" + serverIP
		}
		if !strings.Contains(serverIP[7:], ":") {
			serverIP = serverIP + ":8765"
		}

		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(serverIP + "/api/health")
		if err != nil {
			return fmt.Errorf("impossible de joindre le PC 1 à l'adresse %s. Vérifiez que le PC 1 est allumé et connecté en Admin", serverIP)
		}
		resp.Body.Close()

		cfg := NetworkConfig{
			Mode:      "client",
			ServerURL: serverIP,
			Token:     defaultNetworkToken,
		}
		data, _ := json.MarshalIndent(cfg, "", "  ")
		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return fmt.Errorf("ecriture configuration reseau: %w", err)
		}
		a.networkConfig = cfg
		return nil
	}

	_ = os.Remove(configPath)
	a.networkConfig = NetworkConfig{Mode: "standalone", Token: defaultNetworkToken}
	return nil
}

func getLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}
