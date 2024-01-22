package database

import (
	"context"
	proto "github.com/papireio/go-users-service/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func UpdateUser(client *mongo.Client, req *proto.UpdateUserRequest) (*User, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

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

	if err := collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u); err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return u, nil
}
