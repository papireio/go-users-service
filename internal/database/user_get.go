package database

import (
	"context"
	proto "github.com/papireio/go-users-service/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func GetUser(client *mongo.Client, req *proto.GetUserRequest) (*User, error) {
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := collection.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	return u, nil
}
