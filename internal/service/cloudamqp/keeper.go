package cloudamqp

import (
	"context"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
	"github.com/xudongzhaodev/freeplan-keeper/pkg/cloudamqp"
)

// Keeper implements the service.Keeper interface for CloudAMQP
type Keeper struct {
	client *cloudamqp.Client
}

// NewKeeper creates a new CloudAMQP keeper instance
func NewKeeper(cfg *config.Config) (*Keeper, error) {
	// Skip if CloudAMQP configuration is missing or disabled
	if cfg.CloudAMQP == nil || !cfg.CloudAMQP.Enabled {
		return nil, nil
	}

	client, err := cloudamqp.NewClient(
		cfg.CloudAMQP.URI,
		cfg.CloudAMQP.Queue,
		cfg.Hostname, // Pass the global hostname
	)
	if err != nil {
		return nil, err
	}

	return &Keeper{
		client: client,
	}, nil
}

// Start performs a ping check to CloudAMQP
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
	return "CloudAMQP"
}

func (k *Keeper) Ping() error {
	if k == nil || k.client == nil {
		return nil
	}
	return k.client.Ping()
} 