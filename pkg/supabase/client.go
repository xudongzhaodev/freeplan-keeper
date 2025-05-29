package supabase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Client struct {
	conn            *pgx.Conn
	keepRecordLimit int
	hostname        string // Global hostname from config
}

// NewClient creates a new Supabase client using the PostgreSQL connection
func NewClient(uri string, keepRecordLimit int, hostname string) (*Client, error) {
	// Create the connection
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	return &Client{
		conn:            conn,
		keepRecordLimit: keepRecordLimit,
		hostname:        hostname,
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

	if _, err := c.conn.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create keep_alive_reserved table: %w", err)
	}

	// Insert a new record
	insertSQL := `
	INSERT INTO keep_alive_reserved (ping_source, ping_details)
	VALUES ($1, $2::jsonb)`

	details := fmt.Sprintf(`{"hostname": "%s", "version": "1.0"}`, c.hostname)

	if _, err := c.conn.Exec(ctx, insertSQL, "supabase-keeper", details); err != nil {
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

	if _, err := c.conn.Exec(ctx, cleanupSQL, c.keepRecordLimit); err != nil {
		// Just log the error but don't fail the ping
		fmt.Printf("Warning: failed to cleanup old records: %v\n", err)
	}

	return nil
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close(context.Background())
	}
	return nil
}
