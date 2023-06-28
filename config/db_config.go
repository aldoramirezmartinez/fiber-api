package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func ConnectDB() (*mongo.Database, error) {
	mongoURI, err := GetMongoURI()
	if err != nil {
		fmt.Printf("Error getting MongoDB URI: %s\n", err)
		return nil, err
	}

	dbName := GetDBName()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %s\n", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to ping MongoDB: %s\n", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	db = client.Database(dbName)

	return db, nil
}

func GetDB() *mongo.Database {
	return db
}
