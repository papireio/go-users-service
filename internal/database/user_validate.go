package database

import (
	"context"
	"errors"
	proto "go-users/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ValidateEmail(client *mongo.Client, req *proto.ValidateEmailRequest) error {
	user, err := GetUser(client, &proto.GetUserRequest{Uuid: req.Uuid})
	if err != nil {
		return err
	}

	if user.ValidationToken != req.ValidationToken {
		return errors.New("invalid_token")
	}

	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}
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

	return collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u)
}
