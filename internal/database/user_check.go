package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	proto "goservice/pkg/api/grpc"
	"time"
)

func CheckEmail(client *mongo.Client, req *proto.CheckEmailRequest) bool {
	u := &User{}
	filter := bson.D{{"email", req.Email}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	return err != nil
}
