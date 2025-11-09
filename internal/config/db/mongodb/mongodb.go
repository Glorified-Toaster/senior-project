// Package mongodb provides MongoDB database configuration and connection handling.
package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
	once     sync.Once
)

// MongoConnect : establishes a connection to the MongoDB database.
func MongoConnect(uri string, dbName string) error {
	var err error

	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create a new MongoDB client options
		clientOptions := options.Client().ApplyURI(uri)

		// Connect to MongoDB
		Client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			err = fmt.Errorf("failed to connect to MongoDB: %w", err)
			return

		}

		// Ping the database
		if err = Client.Ping(ctx, nil); err != nil {
			_ = Client.Disconnect(ctx)
			Client = nil
			err = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
		}

		// Get a handle for the database
		Database = Client.Database(dbName)

		log.Printf("Connected to MongoDB! Database : %s", dbName)
	})
	return err
}

// MongoDisconnect : disconnects from the MongoDB database.
func MongoDisconnect() error {
	// check if client is initialized
	if Client == nil {
		log.Println("MongoDB client is not initialized.")
		return nil // nil error since there's nothing to disconnect
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// disconnect
	if err := Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}
	// reset
	Client = nil
	Database = nil

	log.Println("Disconnected from MongoDB!")
	return nil
}

// GetCollection : retrieves a collection from the MongoDB database.
func GetCollection(collectionName string) *mongo.Collection {
	if Database == nil {
		panic("Database is not initialized. Call MongoConnect() first.")
	}
	return Database.Collection(collectionName)
}

// MakeURI : construct MongoDB URI
func MakeURI(host, port, uname, passwd, db string) string {
	if uname != "" && passwd != "" {
		// mongodb://user:pass@localhost:27017/database
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", uname, passwd, host, port, db)
	}
	// mongodb://localhost:27017/database
	return fmt.Sprintf("mongodb://%s:%s/%s", host, port, db)
}
