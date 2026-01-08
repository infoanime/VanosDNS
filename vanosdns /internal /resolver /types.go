package resolver

import (
	"time"

	"github.com/miekg/dns"
)

// Query représente une requête DNS normalisée
type Query struct {
	Name      string
	QType     uint16
	QClass    uint16
	Timestamp time.Time
	ClientIP  string
}

// Response représente une réponse DNS interne
type Response struct {
	Msg       *dns.Msg
	FromCache bool
	Latency   time.Duration
}

// ResolveResult est utilisé entre resolver → router
type ResolveResult struct {
	Response *dns.Msg
	Err      error
	Source   ResolveSource
	Latency  time.Duration
}

// ResolveSource indique d’où vient la réponse
type ResolveSource string

const (
	SourceCache    ResolveSource = "cache"
	SourceMesh     ResolveSource = "mesh"
	SourceUpstream ResolveSource = "upstream"
	SourceUnknown  ResolveSource = "unknown"
)

// ResolverStats sert au monitoring / GUI
type ResolverStats struct {
	TotalQueries   uint64
	CacheHits      uint64
	CacheMisses    uint64
	UpstreamCalls  uint64
	MeshCalls      uint64
	Errors         uint64
	AvgLatencyMS   float64
	P95LatencyMS   float64
	LastUpdateTime time.Time
}
