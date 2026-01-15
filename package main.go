package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set MongoDB URI and connect
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	// Get collection handle
	collection := client.Database("exampledb").Collection("examples")

	// Insert a document
	doc := bson.D{
		bson.E{Key: "name", Value: "MongoDB Go Driver"},
		bson.E{Key: "version", Value: "v1.17.x"},
	}
	insertResult, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	fmt.Printf("Inserted document with _id: %v\n", insertResult.InsertedID)

	// Find the document
	var result bson.M
	err = collection.FindOne(ctx, bson.D{
		bson.E{Key: "name", Value: "MongoDB Go Driver"},
	}).Decode(&result)
	if err != nil {
		log.Fatalf("FindOne failed: %v", err)
	}
	fmt.Printf("Found document: %+v\n", result)
}
