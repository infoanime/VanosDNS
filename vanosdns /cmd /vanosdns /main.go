package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vanosdns/internal/antisibil"
	"vanosdns/internal/cache"
	"vanosdns/internal/config"
	"vanosdns/internal/forwarder"
	"vanosdns/internal/gui"
	"vanosdns/internal/mesh"
	"vanosdns/internal/metrics"
	"vanosdns/internal/resolver"
	"vanosdns/internal/routing"
)

func main() {
	fmt.Println("VanosDNS v1 – starting up")

	// Load configuration
	cfg, err := config.Load("configs/vanosdns.yaml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Init metrics
	metrics.Init()

	// Init cache
	dnsCache := cache.New(cfg.Cache)

	// Init upstream forwarder
	forwarder := forwarder.New(cfg.Upstreams, cfg.Timeout)

	// Init mesh network
	meshNet := mesh.New(cfg.Mesh)

	// Init Anti-Sybil engine
	antiSybil := antisibil.New()

	// Init router (mesh + upstream + hedged requests)
	router := routing.NewRouter(meshNet, forwarder, antiSybil)

	// Start DNS UDP server
	go func() {
		log.Printf("DNS server listening on %s", cfg.Listen)
		if err := resolver.ListenUDP(cfg.Listen, dnsCache, router); err != nil {
			log.Fatalf("DNS server crashed: %v", err)
		}
	}()

	// Start local GUI
	go func() {
		log.Printf("GUI listening on :%d", cfg.GUI.Port)
		http.Handle("/", gui.Handler(dnsCache, meshNet, metrics.Global()))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.GUI.Port), nil))
	}()

	// Background maintenance loop
	go maintenanceLoop(dnsCache, meshNet, antiSybil)

	// Graceful shutdown
	waitForShutdown()
	fmt.Println("VanosDNS stopped cleanly")
}

func maintenanceLoop(
	cache *cache.Cache,
	meshNet *mesh.Mesh,
	antiSybil *antisibil.Engine,
) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cache.PurgeExpired()
		meshNet.RefreshNeighbors()
		antiSybil.Update()
		metrics.Tick()
	}
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
