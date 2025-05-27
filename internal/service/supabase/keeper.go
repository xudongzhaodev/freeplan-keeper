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
func NewKeeper(cfg *config.Config) (*Keeper, error) {
	// Skip if Supabase configuration is missing or disabled
	if cfg.Supabase == nil || !cfg.Supabase.Enabled {
		return nil, nil
	}

	client, err := supabase.NewClient(
		cfg.Supabase.URI,
		cfg.Supabase.KeepRecordLimit,
		cfg.Hostname,  // Pass the global hostname
	)
	if err != nil {
		return nil, err
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

func (k *Keeper) Ping() error {
	if k == nil || k.client == nil {
		return nil
	}
	return k.client.Ping()
} 