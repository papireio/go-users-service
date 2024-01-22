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

func ValidateEmail(client *mongo.Client, req *proto.ValidateEmailRequest) error {
	if req.Uuid == "" || req.ValidationToken == "" {
		return status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

	user, err := GetUser(client, &proto.GetUserRequest{Uuid: req.Uuid})
	if err != nil {
		return status.Error(codes.NotFound, "User not found")
	}

	if user.ValidationToken != req.ValidationToken {
		return status.Error(codes.Unauthenticated, "Token invalid")
	}

	u := &User{}
	filter := bson.D{{"uuid", user.Uuid}}
	payload := bson.M{"$set": bson.M{"validation_token": ""}}

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
		return status.Error(codes.Internal, "Internal server error (getting recently exist user)")
	}

	return nil
}
