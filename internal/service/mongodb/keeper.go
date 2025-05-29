package mongodb

import (
	"context"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
	"github.com/xudongzhaodev/freeplan-keeper/pkg/mongodb"
)

// Keeper implements the service.Keeper interface for MongoDB
type Keeper struct {
	client *mongodb.Client
}

// NewKeeper creates a new MongoDB keeper instance
func NewKeeper(cfg *config.Config) (*Keeper, error) {
	// Skip if MongoDB configuration is missing or disabled
	if cfg.MongoDB == nil || !cfg.MongoDB.Enabled {
		return nil, nil
	}

	client, err := mongodb.NewClient(
		cfg.MongoDB.URI,
		cfg.MongoDB.Database,
		cfg.MongoDB.KeepRecordLimit,
		cfg.Hostname, // Pass the global hostname
	)
	if err != nil {
		return nil, err
	}

	return &Keeper{
		client: client,
	}, nil
}

// Start performs a ping check to MongoDB
func (k *Keeper) Start(ctx context.Context) error {
	return k.client.Ping()
}

// Stop performs cleanup
func (k *Keeper) Stop() error {
	if k.client != nil {
		return k.client.Close()
	}
	return nil
}

// Name returns the service identifier
func (k *Keeper) Name() string {
	return "MongoDB"
}

func (k *Keeper) Ping() error {
	if k == nil || k.client == nil {
		return nil
	}
	return k.client.Ping()
}
