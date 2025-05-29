package service

import (
	"context"
	"log"
)

// Keeper interface defines the methods that each service keeper must implement
type Keeper interface {
	Start(ctx context.Context) error
	Stop() error
	Name() string
}

// Manager handles all keeper instances
type Manager struct {
	keepers []Keeper
}

// NewManager creates a new keeper manager
func NewManager() *Manager {
	return &Manager{
		keepers: make([]Keeper, 0),
	}
}

// RegisterKeeper adds a new keeper to the manager
func (m *Manager) RegisterKeeper(k Keeper) {
	m.keepers = append(m.keepers, k)
}

// RunOnce executes one check for all registered keepers
func (m *Manager) RunOnce(ctx context.Context) error {
	for _, keeper := range m.keepers {
		if err := keeper.Start(ctx); err != nil {
			log.Printf("[%s] Check failed: %v", keeper.Name(), err)
		} else {
			log.Printf("[%s] Check successful", keeper.Name())
		}
	}
	return nil
}

// IsEmpty returns true if no keepers are registered
func (m *Manager) IsEmpty() bool {
	return len(m.keepers) == 0
}

// Stop gracefully shuts down all keepers
func (m *Manager) Stop() {
	for _, keeper := range m.keepers {
		if err := keeper.Stop(); err != nil {
			log.Printf("[%s] Error stopping keeper: %v", keeper.Name(), err)
		}
	}
}
