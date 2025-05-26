package supabase

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool            *pgxpool.Pool
	keepRecordLimit int
}

// NewClient creates a new Supabase client using the PostgreSQL connection
func NewClient(url, apiKey string, keepRecordLimit int) (*Client, error) {
	// Convert Supabase URL to PostgreSQL connection string
	// Format: postgres://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
	// Remove https:// if present
	dbURL := url
	if len(dbURL) > 8 && dbURL[:8] == "https://" {
		dbURL = dbURL[8:]
	}
	connStr := fmt.Sprintf("postgres://postgres:%s@db.%s:5432/postgres", apiKey, dbURL)

	// Create a connection pool configuration
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Set some reasonable pool settings
	config.MaxConns = 1 // We only need one connection for this use case
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{
		pool:            pool,
		keepRecordLimit: keepRecordLimit,
	}, nil
}

// Ping checks database connectivity and maintains activity by inserting a record
func (c *Client) Ping() error {
	ctx := context.Background()

	// Create the keep_alive_reserved table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS keep_alive_reserved (
		id BIGSERIAL PRIMARY KEY,
		ping_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		ping_source VARCHAR(255),
		ping_details JSONB DEFAULT '{}'::jsonb
	)`

	if _, err := c.pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create keep_alive_reserved table: %w", err)
	}

	// Insert a new record
	insertSQL := `
	INSERT INTO keep_alive_reserved (ping_source, ping_details)
	VALUES ($1, $2::jsonb)`

	details := fmt.Sprintf(`{"host": "%s", "version": "1.0"}`, "freeplan-keeper")
	
	if _, err := c.pool.Exec(ctx, insertSQL, "supabase-keeper", details); err != nil {
		return fmt.Errorf("failed to insert keep-alive record: %w", err)
	}

	// Clean up old records based on configured limit
	cleanupSQL := `
	DELETE FROM keep_alive_reserved
	WHERE id NOT IN (
		SELECT id FROM keep_alive_reserved
		ORDER BY ping_timestamp DESC
		LIMIT $1
	)`

	if _, err := c.pool.Exec(ctx, cleanupSQL, c.keepRecordLimit); err != nil {
		// Just log the error but don't fail the ping
		fmt.Printf("Warning: failed to cleanup old records: %v\n", err)
	}

	return nil
}

// Close closes the database connection pool
func (c *Client) Close() error {
	if c.pool != nil {
		c.pool.Close()
	}
	return nil
} 