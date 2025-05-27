package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration settings for the application
type Config struct {
	Hostname string `yaml:"hostname"`  // Global identifier for this keeper instance
	MongoDB struct {
		Enabled         bool   `yaml:"enabled"`
		URI            string `yaml:"uri"`
		Database       string `yaml:"database"`
		KeepRecordLimit int    `yaml:"keep_records_limit"`
	} `yaml:"mongodb"`
	
	Supabase struct {
		Enabled         bool   `yaml:"enabled"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		DBName         string `yaml:"dbname"`
		KeepRecordLimit int    `yaml:"keep_records_limit"`
	} `yaml:"supabase"`
}

// Load reads and parses the configuration file
func Load() (*Config, error) {
	// Try to read from config.yaml first
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		// If config.yaml doesn't exist, try config.yml
		data, err = os.ReadFile("config.yml")
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default values
	if cfg.Hostname == "" {
		cfg.Hostname = "freeplan-keeper" // default hostname identifier
	}
	if cfg.MongoDB.KeepRecordLimit <= 0 {
		cfg.MongoDB.KeepRecordLimit = 100 // default to keeping 100 records
	}
	if cfg.Supabase.KeepRecordLimit <= 0 {
		cfg.Supabase.KeepRecordLimit = 100 // default to keeping 100 records
	}
	if cfg.Supabase.Port <= 0 {
		cfg.Supabase.Port = 5432 // default port for Supabase Pooler
	}

	// Validate required fields for enabled services
	if cfg.MongoDB.Enabled {
		if cfg.MongoDB.URI == "" {
			return nil, fmt.Errorf("mongodb.uri is required when mongodb is enabled")
		}
		if cfg.MongoDB.Database == "" {
			return nil, fmt.Errorf("mongodb.database is required when mongodb is enabled")
		}
	}

	if cfg.Supabase.Enabled {
		if cfg.Supabase.Host == "" {
			return nil, fmt.Errorf("supabase.host is required when supabase is enabled")
		}
		if cfg.Supabase.User == "" {
			return nil, fmt.Errorf("supabase.user is required when supabase is enabled")
		}
		if cfg.Supabase.Password == "" {
			return nil, fmt.Errorf("supabase.password is required when supabase is enabled")
		}
		if cfg.Supabase.DBName == "" {
			return nil, fmt.Errorf("supabase.dbname is required when supabase is enabled")
		}
	}

	return &cfg, nil
} 