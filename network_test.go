package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadNetworkConfig(t *testing.T) {
	t.Run("defaults to standalone", func(t *testing.T) {
		app := &App{dataDir: t.TempDir()}
		if err := app.loadNetworkConfigAt(filepath.Join(app.dataDir, networkConfigFilename)); err != nil {
			t.Fatalf("loadNetworkConfig() error = %v", err)
		}
		if app.networkConfig.Mode != "standalone" {
			t.Fatalf("mode = %q, want standalone", app.networkConfig.Mode)
		}
	})

	t.Run("rejects a short token without enabling server mode", func(t *testing.T) {
		dataDir := t.TempDir()
		config := `{"mode":"server","token":"short"}`
		if err := os.WriteFile(filepath.Join(dataDir, networkConfigFilename), []byte(config), 0600); err != nil {
			t.Fatal(err)
		}
		app := &App{dataDir: dataDir}
		if err := app.loadNetworkConfigAt(filepath.Join(dataDir, networkConfigFilename)); err == nil {
			t.Fatal("loadNetworkConfig() expected a token validation error")
		}
		if app.networkConfig.Mode != "standalone" {
			t.Fatalf("invalid config enabled mode %q", app.networkConfig.Mode)
		}
	})

	t.Run("accepts a valid client config", func(t *testing.T) {
		dataDir := t.TempDir()
		config := `{"mode":"client","serverUrl":"http://192.168.1.20:8765","token":"` + strings.Repeat("a", 32) + `"}`
		if err := os.WriteFile(filepath.Join(dataDir, networkConfigFilename), []byte(config), 0600); err != nil {
			t.Fatal(err)
		}
		app := &App{dataDir: dataDir}
		if err := app.loadNetworkConfigAt(filepath.Join(dataDir, networkConfigFilename)); err != nil {
			t.Fatalf("loadNetworkConfig() error = %v", err)
		}
		if app.networkConfig.Mode != "client" {
			t.Fatalf("mode = %q, want client", app.networkConfig.Mode)
		}
	})
}

func TestNetworkHandlerRequiresTokenAndDispatchesRPC(t *testing.T) {
	const token = "12345678901234567890123456789012"
	handler := newNetworkHandler(&App{networkConfig: NetworkConfig{Token: token}})

	t.Run("rejects missing token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/rpc", strings.NewReader(`{"method":"GetAPropos","args":[]}`))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
		}
	})

	t.Run("returns allowed method result", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/rpc", strings.NewReader(`{"method":"GetAPropos","args":[]}`))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
		}
		var envelope struct {
			Result map[string]string `json:"result"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if envelope.Result["nom"] != "Gestion de Stock" {
			t.Fatalf("result = %#v", envelope.Result)
		}
	})

	t.Run("rejects unknown methods", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/api/rpc", strings.NewReader(`{"method":"DeleteEverything","args":[]}`))
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
		}
	})
}

func TestNetworkSessionsEnforceRolePermissions(t *testing.T) {
	app := &App{
		rpcSessions: map[string]rpcSession{
			"user-session": {
				role:      "user",
				expiresAt: time.Now().Add(time.Hour),
			},
			"admin-session": {
				role:      "admin",
				expiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	request := rpcRequest{
		Method:       "CreerProduit",
		Args:         []json.RawMessage{json.RawMessage(`{}`)},
		SessionToken: "user-session",
	}
	if _, err := app.dispatchRPCRequest(request); err == nil || !strings.Contains(err.Error(), "accès refusé") {
		t.Fatalf("user product creation error = %v, want role denial", err)
	}

	request.Method = "ListerMouvementsStock"
	request.Args = []json.RawMessage{json.RawMessage(`{}`)}
	if _, err := app.dispatchRPCRequest(request); err == nil || !strings.Contains(err.Error(), "accès refusé") {
		t.Fatalf("user movement listing error = %v, want role denial", err)
	}

	request.Method = "ListerFactures"
	if _, err := app.dispatchRPCRequest(request); err == nil || !strings.Contains(err.Error(), "base de données") {
		t.Fatalf("user invoice listing error = %v, want authorization success before database check", err)
	}

	request.Method = "CreerProduit"
	request.Args = []json.RawMessage{json.RawMessage(`{}`)}
	request.SessionToken = "admin-session"
	if _, err := app.dispatchRPCRequest(request); err == nil || !strings.Contains(err.Error(), "base de données") {
		t.Fatalf("admin product creation error = %v, want database-unavailable after authorization", err)
	}

	request.SessionToken = "expired-session"
	if _, err := app.dispatchRPCRequest(request); err == nil || !strings.Contains(err.Error(), "session invalide") {
		t.Fatalf("invalid session error = %v, want session denial", err)
	}
}

func TestAdminAutoStartPartageServer(t *testing.T) {
	app := &App{
		dataDir: t.TempDir(),
	}

	info := app.GetServerInfo()
	if info.Running {
		t.Fatal("expected server to not be running initially")
	}

	app.startPartageServer()
	infoAfter := app.GetServerInfo()
	app.stopPartageServer()

	if infoAfter.Port != 8765 {
		t.Fatalf("expected port 8765, got %d", infoAfter.Port)
	}
}
