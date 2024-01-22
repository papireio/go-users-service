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

func FlushSessions(client *mongo.Client, req *proto.FlushSessionsRequest) error {
	if req.SessionToken == "" || req.Uuid == "" {
		return status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

	u := &User{}
	filter := bson.D{{"uuid", req.Uuid}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return status.Error(codes.NotFound, "User not found")
	}

	session := Session{}
	for _, v := range u.Sessions {
		if v.Token == req.SessionToken {
			session = v
		}
	}

	sessions := make([]Session, 0)
	if session.Token != "" {
		sessions = append(sessions, session)
	}

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

	if err := collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u); err != nil {
		return status.Error(codes.Internal, "Internal server error (getting recently exist user)")
	}

	return nil
}
