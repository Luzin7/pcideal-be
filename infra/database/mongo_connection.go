package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	serverSelectionTimeout = 10 * time.Second
	maxPoolSize            = 20
	connectTimeout         = 30 * time.Second
)

func MongoConnection(connectionString string, databaseName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetServerSelectionTimeout(serverSelectionTimeout).
		SetMaxPoolSize(maxPoolSize).
		SetRetryWrites(true)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(context.Background())
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return client.Database(databaseName), nil
}
