package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	proto "goservice/pkg/api/grpc"
	"time"
)

func remove(s []Session, t string) []Session {
	for i, v := range s {
		if v.Token == t {
			return append(s[:i], s[i+1:]...)
		}
	}

	return s
}

func DeleteSession(client *mongo.Client, req *proto.DeleteSessionRequest) error {
	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return err
	}

	sessions := remove(u.Sessions, req.SessionToken)

	filter = bson.D{{"uuid", u.Uuid}}
	payload := bson.M{"$set": bson.M{"sessions": sessions}}

	upsert := false
	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	collection = client.Database("papireio").Collection("users")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u)
}
