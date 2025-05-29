package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	client          *mongo.Client
	db              *mongo.Database
	keepRecordLimit int
	hostname        string // Global hostname from config
}

// NewClient creates a new MongoDB client
func NewClient(uri string, database string, keepRecordLimit int, hostname string) (*Client, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	ctx := context.Background()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Get database instance
	db := client.Database(database)

	return &Client{
		client:          client,
		db:              db,
		keepRecordLimit: keepRecordLimit,
		hostname:        hostname,
	}, nil
}

// Ping checks database connectivity and maintains activity by inserting a record
func (c *Client) Ping() error {
	ctx := context.Background()

	// Get or create the keep_alive_reserved collection
	collection := c.db.Collection("keep_alive_reserved")

	// Create indexes if they don't exist
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "ping_timestamp", Value: -1}}, // Descending index on ping_timestamp
	})
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Drop the problematic id index if it exists
	_, err = collection.Indexes().DropOne(ctx, "id_1")
	if err != nil {
		// Ignore if index doesn't exist, log other errors
		if !strings.Contains(err.Error(), "index not found") {
			fmt.Printf("Warning: failed to drop id index: %v\n", err)
		}
	}

	// Insert a new record
	record := bson.D{
		{Key: "_id", Value: primitive.NewObjectID()},
		{Key: "ping_timestamp", Value: time.Now()},
		{Key: "ping_source", Value: "mongodb-keeper"},
		{Key: "ping_details", Value: bson.D{
			{Key: "hostname", Value: c.hostname},
			{Key: "version", Value: "1.0"},
		}},
	}

	if _, err := collection.InsertOne(ctx, record); err != nil {
		return fmt.Errorf("failed to insert keep-alive record: %w", err)
	}

	// Clean up old records based on configured limit
	// First, find the timestamp of the Nth record (where N is the keep_record_limit)
	var cutoffDoc struct {
		PingTimestamp time.Time `bson:"ping_timestamp"`
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "ping_timestamp", Value: -1}}).SetSkip(int64(c.keepRecordLimit))
	err = collection.FindOne(ctx, bson.M{}, opts).Decode(&cutoffDoc)

	if err != nil && err != mongo.ErrNoDocuments {
		// Log the error but don't fail the ping
		fmt.Printf("Warning: failed to find cutoff timestamp: %v\n", err)
		return nil
	}

	if err != mongo.ErrNoDocuments {
		// Delete records older than the cutoff timestamp
		_, err = collection.DeleteMany(ctx, bson.M{
			"ping_timestamp": bson.M{"$lt": cutoffDoc.PingTimestamp},
		})
		if err != nil {
			// Log the error but don't fail the ping
			fmt.Printf("Warning: failed to cleanup old records: %v\n", err)
		}
	}

	return nil
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Disconnect(context.Background())
	}
	return nil
}
