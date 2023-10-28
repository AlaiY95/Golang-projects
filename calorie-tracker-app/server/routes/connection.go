package routes

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBinstance creates and returns a MongoDB client.
func DBinstance() *mongo.Client {
	// Define the MongoDB connection URI.
	MongoDb := "mongodb://localhost:27017/caloriesdb"

	// Create a new MongoDB client using the provided connection URI.
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	// Create a context with a timeout to connect to the database.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to connect to the MongoDB server.
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Print a message indicating a successful MongoDB connection.
	fmt.Println("Connected to MongoDB")

	return client
}

// Client is a MongoDB client instance created during initialization.
var Client *mongo.Client = DBinstance()

// OpenCollection opens a specific MongoDB collection using the provided client and collection name.
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Create a reference to the specified collection within the "caloriesdb" database.
	var collection *mongo.Collection = client.Database("caloriesdb").Collection(collectionName)
	return collection
}
