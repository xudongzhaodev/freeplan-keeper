package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xzhao/freeplan-keeper/internal/config"
	"github.com/xzhao/freeplan-keeper/pkg/supabase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Supabase client if enabled
	if cfg.Supabase.Enabled {
		supabaseClient, err := supabase.NewClient(cfg.Supabase.URL, cfg.Supabase.DBPassword, cfg.Supabase.KeepRecordLimit)
		if err != nil {
			log.Fatalf("Failed to initialize Supabase client: %v", err)
		}
		defer supabaseClient.Close()

		if err := supabaseClient.Ping(); err != nil {
			log.Printf("Supabase ping failed: %v", err)
		} else {
			log.Println("Supabase ping successful")
		}
	}

	log.Println("All services checked successfully")
} 