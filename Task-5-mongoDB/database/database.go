package database

import (
	"context"
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDataBase(uri, dbName string) (*mongo.Database, error) {

	setClient := options.Client().ApplyURI(uri)

	// Creaing a new MongoDB client
	client, err := mongo.Connect(context.TODO(), setClient)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Set a timeout for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client.Database(dbName), nil
}
