package database

import (
	"context"
	"github.com/papireio/go-users-service/internal/utils"
	proto "github.com/papireio/go-users-service/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func CreateSession(client *mongo.Client, req *proto.CreateSessionRequest) (*User, string, error) {
	if req.Email == "" || req.Password == "" {
		return nil, "", status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

	u := &User{}
	filter := bson.D{{"email", req.Email}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return nil, "", status.Error(codes.NotFound, "User not found")
	}

	if isPasswordValid := utils.CompareHashPassword(req.Password, u.PasswordHash); !isPasswordValid {
		return nil, "", status.Error(codes.PermissionDenied, "Password mismatch")
	}

	sessionToken, err := utils.GetToken()
	if err != nil {
		return nil, "", status.Error(codes.Internal, "Internal server error (getting session token)")
	}

	sessions := append(u.Sessions, Session{
		Token:     sessionToken,
		CreatedAt: time.Now().String(),
	})

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

	if err = collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u); err != nil {
		return nil, "", status.Error(codes.Internal, "Internal server error (getting recently exist user)")
	}

	return u, sessionToken, nil
}
