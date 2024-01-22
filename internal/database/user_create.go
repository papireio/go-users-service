package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/papireio/go-users-service/internal/utils"
	proto "github.com/papireio/go-users-service/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func CreateUser(client *mongo.Client, req *proto.CreateUserRequest) (*User, string, error) {
	if req.Email == "" || req.Password == "" {
		return nil, "", status.Error(codes.InvalidArgument, "Incorrect request argument")
	}

	u := &User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return u, "", err
	}

	sessionToken, err := utils.GetToken()
	if err != nil {
		return u, "", err
	}

	validationToken, err := utils.GetToken()
	if err != nil {
		return u, "", err
	}

	sessions := []Session{
		{
			Token:     sessionToken,
			CreatedAt: time.Now().String(),
		},
	}

	collection := client.Database("papireio").Collection("users")
	res, err := collection.InsertOne(ctx, bson.D{
		{"email", req.Email},
		{"hash", hash},
		{"uuid", uuid.New().String()},
		{"sessions", sessions},
		{"validation_token", validationToken},
	})
	if err != nil {
		return nil, "", err
	}

	id := res.InsertedID
	filter := bson.D{{"_id", id}}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&u)
	// TODO: Send validation email

	return u, sessionToken, err
}
