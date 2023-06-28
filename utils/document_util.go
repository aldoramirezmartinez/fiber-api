package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckDocumentExists(collection *mongo.Collection, documentID primitive.ObjectID) (bool, error) {
	ctx := context.TODO()

	count, err := collection.CountDocuments(ctx, bson.M{"_id": documentID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
