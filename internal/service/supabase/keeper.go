package supabase

import (
	"context"
	"fmt"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
	"github.com/xudongzhaodev/freeplan-keeper/pkg/supabase"
)

// Keeper implements the service.Keeper interface for Supabase
type Keeper struct {
	client *supabase.Client
}

// NewKeeper creates a new Supabase keeper instance
func NewKeeper(cfg *config.Config) (*Keeper, error) {
	client, err := supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}
	return &Keeper{
		client: client,
	}, nil
}

// Start performs a ping check to Supabase
func (k *Keeper) Start(ctx context.Context) error {
	return k.client.Ping()
}

// Stop performs cleanup (no-op for Supabase)
func (k *Keeper) Stop() error {
	if k.client != nil {
		return k.client.Close()
	}
	return nil
}

// Name returns the service identifier
func (k *Keeper) Name() string {
	return "Supabase"
} 