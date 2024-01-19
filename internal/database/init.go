package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func Init(client *mongo.Client) error {
	db := client.Database("papireio")
	if err := db.CreateCollection(ctx, "users"); err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "uuid", Value: 2}},
			Options: options.Index().SetUnique(true),
		},
	}

	collection := db.Collection("users")
	if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}

	return nil
}
