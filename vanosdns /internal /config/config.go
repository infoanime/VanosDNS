package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// Config struct principale VanosDNS
type Config struct {
	Listen   string       `yaml:"listen"`   // ex: ":53"
	GUI      GUIConfig    `yaml:"gui"`
	Cache    CacheConfig  `yaml:"cache"`
	Mesh     MeshConfig   `yaml:"mesh"`
	Upstreams []string    `yaml:"upstreams"` // ex: ["1.1.1.1:53","8.8.8.8:53"]
	TimeoutMS int         `yaml:"timeout_ms"`
}

// GUI Config pour le serveur de stats local
type GUIConfig struct {
	Port int `yaml:"port"`
}

// Cache Config
type CacheConfig struct {
	DefaultTTLSeconds int `yaml:"default_ttl_seconds"`
	MaxEntries        int `yaml:"max_entries"`
}

// Mesh Config
type MeshConfig struct {
	DiscoveryCount int `yaml:"discovery_count"` // nbr de noeuds proches connus
}

// Load lit et parse le YAML
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	setDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Défauts si non définis
func setDefaults(cfg *Config) {
	if cfg.Listen == "" {
		cfg.Listen = ":53"
	}
	if cfg.GUI.Port == 0 {
		cfg.GUI.Port = 8080
	}
	if cfg.Cache.DefaultTTLSeconds <= 0 {
		cfg.Cache.DefaultTTLSeconds = 30
	}
	if cfg.Cache.MaxEntries <= 0 {
		cfg.Cache.MaxEntries = 10000
	}
	if cfg.Mesh.DiscoveryCount <= 0 {
		cfg.Mesh.DiscoveryCount = 4
	}
	if cfg.TimeoutMS <= 0 {
		cfg.TimeoutMS = 800
	}
}

// Validate vérifie que la config est cohérente
func validate(cfg *Config) error {
	if len(cfg.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream must be defined")
	}
	if cfg.Cache.DefaultTTLSeconds < 1 {
		return fmt.Errorf("cache TTL must be positive")
	}
	if cfg.GUI.Port <= 0 || cfg.GUI.Port > 65535 {
		return fmt.Errorf("invalid GUI port")
	}
	if cfg.Mesh.DiscoveryCount < 1 {
		return fmt.Errorf("mesh.discovery_count must be >=1")
	}
	if cfg.TimeoutMS < 100 {
		return fmt.Errorf("timeout_ms too low (<100ms)")
	}
	return nil
}
