package bootstrap

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDatabase creates a new MongoDB database connection using the provided environment configuration.
// It takes an `Env` struct pointer as a parameter and returns a `mongo.Client` instance.
// The `Env` struct contains the necessary information for connecting to the MongoDB database, such as the host and port.
// The function establishes a connection to the MongoDB server using the provided host and port,
// and returns the client instance for further use.
// If any error occurs during the connection establishment or ping, the function will log the error and terminate the program.
func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort

	mongodbURI := "mongodb://" + dbHost + ":" + dbPort

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return *client
}

// CloseMongoDatabase closes the connection to the MongoDB database.
// It takes a *mongo.Client as a parameter and disconnects the client from the database.
// If the client is nil, it returns immediately without performing any action.
// If an error occurs during disconnection, it logs the error and exits the program.
// After successful disconnection, it logs a message indicating that the connection to MongoDB has been closed.
func CloseMongoDatabase(client *mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
