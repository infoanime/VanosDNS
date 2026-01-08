package metrics

import (
	"sync"
	"time"
)

// Stats contient les métriques globales du resolver
type Stats struct {
	mu            sync.RWMutex
	TotalQueries  uint64
	CacheHits     uint64
	CacheMisses   uint64
	UpstreamCalls uint64
	MeshCalls     uint64
	Errors        uint64
	TotalLatency  time.Duration
}

// singleton
var globalStats *Stats
var once sync.Once

// Init initialise le singleton Metrics
func Init() {
	once.Do(func() {
		globalStats = &Stats{}
	})
}

// Global retourne le singleton metrics
func Global() *Stats {
	return globalStats
}

// Tick à appeler périodiquement pour update les stats
func Tick() {
	globalStats.mu.Lock()
	defer globalStats.mu.Unlock()
	// futur : calcul P95/P99, logging, export Prometheus…
}

// Increment helpers
func IncTotalQueries() {
	globalStats.mu.Lock()
	globalStats.TotalQueries++
	globalStats.mu.Unlock()
}

func IncCacheHit() {
	globalStats.mu.Lock()
	globalStats.CacheHits++
	globalStats.mu.Unlock()
}

func IncCacheMiss() {
	globalStats.mu.Lock()
	globalStats.CacheMisses++
	globalStats.mu.Unlock()
}

func IncUpstream() {
	globalStats.mu.Lock()
	globalStats.UpstreamCalls++
	globalStats.mu.Unlock()
}

func IncMesh() {
	globalStats.mu.Lock()
	globalStats.MeshCalls++
	globalStats.mu.Unlock()
}

func IncError() {
	globalStats.mu.Lock()
	globalStats.Errors++
	globalStats.mu.Unlock()
}

// AddLatency accumule la latence pour calcul d’avg
func AddLatency(d time.Duration) {
	globalStats.mu.Lock()
	globalStats.TotalLatency += d
	globalStats.mu.Unlock()
}

// AverageLatency retourne la latence moyenne en ms
func AverageLatency() float64 {
	globalStats.mu.RLock()
	defer globalStats.mu.RUnlock()
	if globalStats.TotalQueries == 0 {
		return 0
	}
	return float64(globalStats.TotalLatency.Milliseconds()) / float64(globalStats.TotalQueries)
}
