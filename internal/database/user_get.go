package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	proto "goservice/pkg/api/grpc"
	"time"
)

func GetUser(client *mongo.Client, req *proto.GetUserRequest) (*User, error) {
	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	return u, err
}
