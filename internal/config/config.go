package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// MongoDBConfig holds MongoDB specific configuration
type MongoDBConfig struct {
	Enabled         bool   `yaml:"enabled"`
	URI             string `yaml:"uri"`
	Database        string `yaml:"database"`
	KeepRecordLimit int    `yaml:"keep_records_limit"`
}

// SupabaseConfig holds Supabase specific configuration
type SupabaseConfig struct {
	Enabled         bool   `yaml:"enabled"`
	URI             string `yaml:"uri"`
	KeepRecordLimit int    `yaml:"keep_records_limit"`
}

// Config holds all configuration settings for the application
type Config struct {
	Hostname string          `yaml:"hostname"`           // Global identifier for this keeper instance
	MongoDB  *MongoDBConfig  `yaml:"mongodb,omitempty"`  // Optional MongoDB configuration
	Supabase *SupabaseConfig `yaml:"supabase,omitempty"` // Optional Supabase configuration
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

	// Initialize MongoDB config with defaults if it exists
	if cfg.MongoDB != nil {
		if cfg.MongoDB.KeepRecordLimit <= 0 {
			cfg.MongoDB.KeepRecordLimit = 100
		}
		// Validate required fields only if MongoDB is enabled
		if cfg.MongoDB.Enabled {
			if cfg.MongoDB.URI == "" {
				return nil, fmt.Errorf("mongodb.uri is required when mongodb is enabled")
			}
			if cfg.MongoDB.Database == "" {
				return nil, fmt.Errorf("mongodb.database is required when mongodb is enabled")
			}
		}
	}

	// Initialize Supabase config with defaults if it exists
	if cfg.Supabase != nil {
		if cfg.Supabase.KeepRecordLimit <= 0 {
			cfg.Supabase.KeepRecordLimit = 100
		}
		// Validate required fields only if Supabase is enabled
		if cfg.Supabase.Enabled {
			if cfg.Supabase.URI == "" {
				return nil, fmt.Errorf("supabase.uri is required when supabase is enabled")
			}
		}
	}

	return &cfg, nil
}
