package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"exellence-materiaux/internal/db"
	"exellence-materiaux/partage/src/controllers"
	"exellence-materiaux/partage/src/routers"
	"exellence-materiaux/partage/src/services"
)

func main() {
	portFlag := flag.Int("port", 8765, "Port d'écoute du serveur")
	dbPathFlag := flag.String("db", "", "Chemin vers la base exellence.db")
	flag.Parse()

	dbPath := *dbPathFlag
	projectDir := ""

	if dbPath == "" {
		candidates := []string{
			"exellence.db",
			filepath.Join("..", "exellence.db"),
			filepath.Join(".", "exellence.db"),
		}
		for _, cand := range candidates {
			if _, err := os.Stat(cand); err == nil {
				dbPath = cand
				break
			}
		}
		if dbPath == "" {
			dbPath = "exellence.db"
		}
	}

	absDBPath, err := filepath.Abs(dbPath)
	if err == nil {
		projectDir = filepath.Dir(absDBPath)
	} else {
		projectDir = "."
	}

	database, err := db.Open(absDBPath)
	if err != nil {
		log.Fatalf("Erreur ouverture base de données %s : %v", absDBPath, err)
	}
	defer database.Close()

	serviceManager := services.NewServiceManager(database.DB, projectDir)
	controllerManager := controllers.NewControllerManager(serviceManager)
	router := routers.NewRouter(controllerManager)

	addr := fmt.Sprintf("0.0.0.0:%d", *portFlag)
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ips := getLocalIPs()

	fmt.Println("============================================================")
	fmt.Println("         EXCELLENCE MATÉRIAUX - SERVEUR DE PARTAGE          ")
	fmt.Println("============================================================")
	fmt.Printf("Base de données : %s\n", absDBPath)
	fmt.Printf("Port d'écoute   : %d\n", *portFlag)
	fmt.Println("Adresses pour connecter le 2ème poste :")
	for _, ip := range ips {
		fmt.Printf("  -> http://%s:%d\n", ip, *portFlag)
	}
	fmt.Printf("  -> http://localhost:%d\n", *portFlag)
	fmt.Println("------------------------------------------------------------")
	fmt.Println("Endpoints API REST disponibles :")
	fmt.Println("  [GET]  /api/health")
	fmt.Println("  [POST] /api/auth/login")
	fmt.Println("  [GET]  /api/produits")
	fmt.Println("  [GET]  /api/factures")
	fmt.Println("  [GET]  /api/bons")
	fmt.Println("  [GET]  /api/stock/alertes")
	fmt.Println("  [GET]  /api/parametres")
	fmt.Println("============================================================")
	fmt.Println("Serveur actif et prêt ! (Appuyez sur Ctrl+C pour arrêter)")
	fmt.Println()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur HTTP : %v", err)
		}
	}()

	<-stopChan
	fmt.Println("\nArrêt du serveur en cours...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Erreur lors de l'arrêt : %v", err)
	}
	fmt.Println("Serveur arrêté proprement.")
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
