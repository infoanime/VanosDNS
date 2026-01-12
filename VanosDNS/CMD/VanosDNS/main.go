package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"vanosdns/internal/p2p"
	"vanosdns/internal/resolver"
	"vanosdns/internal/scoring"
)

func main() {
	log.Println("=== VanosDNS Node Starting ===")

	// On crée le dossier data s'il n'existe pas pour éviter les erreurs d'écriture
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Impossible de créer le dossier data : %v", err)
	}

	// Si data/node.priv existe, il le charge, sinon il le crée.
	nodeIdentity, err := p2p.LoadOrGenerateIdentity("data/node.priv")
	if err != nil {
		log.Fatalf("Erreur Identité : %v", err)
	}
	log.Printf("[ID] NodeID chargé : %s", nodeIdentity.NodeID)

	// Gère la persistance des requêtes entre les redémarrages
	cacheL2, err := resolver.NewSQLiteCache("data/cache.db")
	if err != nil {
		log.Fatalf("Erreur Cache L2 : %v", err)
	}
	defer cacheL2.Close()
	log.Println("[CACHE] Base SQLite opérationnelle")

	// Analyse la fiabilité des voisins en temps réel
	scoreEngine := scoring.NewEngine()

	// S'occupe du TCP Hole Punching et du Discovery
	p2pManager := p2p.NewManager(nodeIdentity, scoreEngine)
	
	// On lance le discovery dans une goroutine pour ne pas bloquer le démarrage
	go func() {
		log.Println("[P2P] Démarrage du Discovery et du Peering...")
		p2pManager.StartDiscovery()
	}()

	// Écoute les requêtes DNS standard
	dnsResolver := resolver.New(cacheL2, p2pManager, scoreEngine)
	
	go func() {
		// Port 53 par défaut. Note: nécessite les droits root sur Linux
		address := ":53"
		log.Printf("[DNS] Serveur en écoute sur %s", address)
		if err := dnsResolver.Start(address); err != nil {
			log.Printf("[FATAL] Erreur serveur DNS : %v", err)
			log.Println("ASTUCE : Vérifiez que vous avez les droits root pour le port 53.")
		}
	}()

	// Attente du signal Ctrl+C ou Kill pour fermer les connexions proprement
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("\n[STOP] Arrêt gracieux du nœud...")
