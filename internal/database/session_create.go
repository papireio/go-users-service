package database

import (
	"context"
	"errors"
	"go-users/internal/utils"
	proto "go-users/pkg/api/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func CreateSession(client *mongo.Client, req *proto.CreateSessionRequest) (string, error) {
	u := &User{}
	filter := bson.D{{"email", req.Email}}

	collection := client.Database("papireio").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		return "", err
	}

	if isPasswordValid := utils.CompareHashPassword(req.Password, u.PasswordHash); !isPasswordValid {
		return "", errors.New("invalid_password")
	}

	sessionToken, err := utils.GetToken()
	if err != nil {
		return "", err
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

	err = collection.FindOneAndUpdate(ctx, filter, payload, opt).Decode(&u)

	return sessionToken, err
}
