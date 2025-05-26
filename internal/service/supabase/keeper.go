package supabase

import (
	"context"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
	"github.com/xudongzhaodev/freeplan-keeper/pkg/supabase"
)

// Keeper implements the service.Keeper interface for Supabase
type Keeper struct {
	client *supabase.Client
}

// NewKeeper creates a new Supabase keeper instance
func NewKeeper(cfg *config.Config) *Keeper {
	return &Keeper{
		client: supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.APIKey),
	}
}

// Start performs a ping check to Supabase
func (k *Keeper) Start(ctx context.Context) error {
	return k.client.Ping()
}

// Stop performs cleanup (no-op for Supabase)
func (k *Keeper) Stop() error {
	return nil // Supabase client doesn't need cleanup
}

// Name returns the service identifier
func (k *Keeper) Name() string {
	return "Supabase"
} 