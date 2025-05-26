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
	if cfg.MongoDB.Enabled {
		mongoKeeper, err := mongodb.NewKeeper(cfg)
		if err != nil {
			log.Fatalf("Failed to create MongoDB keeper: %v", err)
		}
		manager.RegisterKeeper(mongoKeeper)
	}

	// Initialize Supabase keeper if enabled
	if cfg.Supabase.Enabled {
		supabaseKeeper := supabase.NewKeeper(cfg)
		manager.RegisterKeeper(supabaseKeeper)
	}

	// Run checks once
	if err := manager.RunOnce(ctx); err != nil {
		log.Printf("Error during checks: %v", err)
	}

	// Cleanup
	manager.Stop()
} 