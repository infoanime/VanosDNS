package cache

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Entry représente une entrée de cache DNS
type Entry struct {
	Msg       *dns.Msg
	ExpiresAt time.Time
}

// Cache est un cache DNS thread-safe avec TTL
type Cache struct {
	mu       sync.RWMutex
	store    map[string]*Entry
	defaultTTL time.Duration
}

// Config permet d'initialiser le cache depuis la config
type Config struct {
	DefaultTTLSeconds int
	MaxEntries        int // futur usage (eviction LRU v1.5)
}

// New initialise un nouveau cache DNS
func New(cfg Config) *Cache {
	ttl := time.Duration(cfg.DefaultTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 30 * time.Second
	}

	return &Cache{
		store:      make(map[string]*Entry),
		defaultTTL: ttl,
	}
}

// Key génère une clé unique pour une requête DNS
func (c *Cache) Key(name string, qtype uint16) string {
	return name + ":" + dns.TypeToString[qtype]
}

// Get récupère une entrée du cache si valide
func (c *Cache) Get(key string) (*dns.Msg, bool) {
	c.mu.RLock()
	entry, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	// TTL expiré → delete
	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		return nil, false
	}

	// Clone pour éviter mutation concurrente
	return entry.Msg.Copy(), true
}

// Set ajoute une réponse DNS au cache
func (c *Cache) Set(key string, msg *dns.Msg) {
	ttl := c.extractTTL(msg)
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	entry := &Entry{
		Msg:       msg.Copy(),
		ExpiresAt: time.Now().Add(ttl),
	}

	c.mu.Lock()
	c.store[key] = entry
	c.mu.Unlock()
}

// PurgeExpired supprime les entrées expirées
func (c *Cache) PurgeExpired() {
	now := time.Now()

	c.mu.Lock()
	for key, entry := range c.store {
		if now.After(entry.ExpiresAt) {
			delete(c.store, key)
		}
	}
	c.mu.Unlock()
}

// Size retourne le nombre d'entrées en cache
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// extractTTL récupère le plus petit TTL des réponses
func (c *Cache) extractTTL(msg *dns.Msg) time.Duration {
	minTTL := uint32(0)
