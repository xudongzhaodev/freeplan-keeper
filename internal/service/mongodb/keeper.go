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
	client, err := mongodb.NewClient(cfg.MongoDB.URI, cfg.MongoDB.Database)
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

// Stop closes the MongoDB connection
func (k *Keeper) Stop() error {
	return k.client.Close()
}

// Name returns the service identifier
func (k *Keeper) Name() string {
	return "MongoDB"
} 