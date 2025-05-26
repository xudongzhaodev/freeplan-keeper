package service

import (
	"context"
	"log"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
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

// Stop gracefully shuts down all keepers
func (m *Manager) Stop() {
	for _, keeper := range m.keepers {
		if err := keeper.Stop(); err != nil {
			log.Printf("[%s] Error stopping keeper: %v", keeper.Name(), err)
		}
	}
}

type Keeper struct {
	cfg      *config.Config
	mongodb  *mongodb.Client
	supabase *supabase.Client
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewKeeper(cfg *config.Config) *Keeper {
	return &Keeper{
		cfg:      cfg,
		stopChan: make(chan struct{}),
	}
}

func (k *Keeper) Start(ctx context.Context) error {
	// Initialize MongoDB client
	mongoClient, err := mongodb.NewClient(k.cfg.MongoDB.URI, k.cfg.MongoDB.Database)
	if err != nil {
		return err
	}
	k.mongodb = mongoClient

	// Initialize Supabase client
	supabaseClient := supabase.NewClient(k.cfg.Supabase.URL, k.cfg.Supabase.APIKey)
	k.supabase = supabaseClient

	// Start the heartbeat routines
	k.wg.Add(2)
	go k.mongodbHeartbeat()
	go k.supabaseHeartbeat()

	return nil
}

func (k *Keeper) Stop() {
	close(k.stopChan)
	k.wg.Wait()
	
	if k.mongodb != nil {
		k.mongodb.Close()
	}
}

func (k *Keeper) mongodbHeartbeat() {
	defer k.wg.Done()
	ticker := time.NewTicker(time.Duration(k.cfg.CheckInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-k.stopChan:
			return
		case <-ticker.C:
			if err := k.mongodb.Ping(); err != nil {
				log.Printf("MongoDB heartbeat failed: %v", err)
			} else {
				log.Println("MongoDB heartbeat successful")
			}
		}
	}
}

func (k *Keeper) supabaseHeartbeat() {
	defer k.wg.Done()
	ticker := time.NewTicker(time.Duration(k.cfg.CheckInterval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-k.stopChan:
			return
		case <-ticker.C:
			if err := k.supabase.Ping(); err != nil {
				log.Printf("Supabase heartbeat failed: %v", err)
			} else {
				log.Println("Supabase heartbeat successful")
			}
		}
	}
} 