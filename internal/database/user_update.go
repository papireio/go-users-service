package database

import (
	"context"
	proto "github.com/papireio/go-users-service/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func UpdateUser(client *mongo.Client, req *proto.UpdateUserRequest) (*User, error) {
	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}
	payload := bson.M{"$set": bson.M{"name": req.Name}}

	upsert := false
	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u)
	return u, err
}
