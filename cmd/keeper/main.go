package main

import (
	"context"
	"log"

	"github.com/xudongzhaodev/freeplan-keeper/internal/config"
	"github.com/xudongzhaodev/freeplan-keeper/internal/service"
	"github.com/xudongzhaodev/freeplan-keeper/internal/service/mongodb"
	"github.com/xudongzhaodev/freeplan-keeper/internal/service/supabase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	manager := service.NewManager()

	// Initialize MongoDB keeper if enabled
	if cfg.MongoDB != nil && cfg.MongoDB.Enabled {
		mongoKeeper, err := mongodb.NewKeeper(cfg)
		if err != nil {
			log.Printf("Warning: Failed to create MongoDB keeper: %v", err)
		} else if mongoKeeper != nil {
			manager.RegisterKeeper(mongoKeeper)
			log.Printf("MongoDB keeper registered successfully")
		}
	}

	// Initialize Supabase keeper if enabled
	if cfg.Supabase != nil && cfg.Supabase.Enabled {
		supabaseKeeper, err := supabase.NewKeeper(cfg)
		if err != nil {
			log.Printf("Warning: Failed to create Supabase keeper: %v", err)
		} else if supabaseKeeper != nil {
			manager.RegisterKeeper(supabaseKeeper)
			log.Printf("Supabase keeper registered successfully")
		}
	}

	// Check if any keepers were registered
	if manager.IsEmpty() {
		log.Printf("Warning: No keepers were successfully initialized")
		return
	}

	// Run checks once
	if err := manager.RunOnce(ctx); err != nil {
		log.Printf("Error during checks: %v", err)
	}

	// Cleanup
	manager.Stop()
}
