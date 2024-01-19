package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"goservice/internal/database"
	"goservice/internal/env"
	proto "goservice/pkg/api/grpc"
	"net"
	"time"
)

type clients struct {
	Mongo *mongo.Client
}

type instance struct {
	proto.UnimplementedGoUsersServer
	clients *clients
}

func Serve(config *env.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", config.Port))
	if err != nil {
		return err
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		return err
	}

	if err = database.Init(mongoClient); err != nil {
		return err
	}

	srv := &instance{clients: &clients{
		Mongo: mongoClient,
	}}

	grpcServer := grpc.NewServer()
	proto.RegisterGoUsersServer(grpcServer, srv)

	return grpcServer.Serve(l)
}

func (i *instance) CreateUser(_ context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user, sessionToken, err := database.CreateUser(i.clients.Mongo, req)
	if err != nil {
		return &proto.CreateUserResponse{
			Success: false,
		}, err
	}

	return &proto.CreateUserResponse{
		Name:         user.Name,
		Email:        user.Email,
		Uuid:         user.Uuid,
		SessionToken: sessionToken,
		Verified:     user.ValidationToken == "",
		Success:      true,
	}, nil
}

func (i *instance) GetUser(_ context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := database.GetUser(i.clients.Mongo, req)
	if err != nil {
		return &proto.GetUserResponse{}, err
	}

	return &proto.GetUserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Uuid:     user.Uuid,
		Verified: user.ValidationToken == "",
	}, nil
}

func (i *instance) UpdateUser(_ context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	user, err := database.UpdateUser(i.clients.Mongo, req)
	if err != nil {
		return &proto.UpdateUserResponse{
			Success: false,
		}, err
	}

	return &proto.UpdateUserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Uuid:     user.Uuid,
		Verified: user.ValidationToken == "",
		Success:  true,
	}, nil
}

func (i *instance) CreateSession(_ context.Context, req *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	token, err := database.CreateSession(i.clients.Mongo, req)
	if err != nil {
		return &proto.CreateSessionResponse{
			Success: false,
		}, err
	}

	return &proto.CreateSessionResponse{
		SessionToken: token,
		Success:      true,
	}, nil
}

func (i *instance) DeleteSession(_ context.Context, req *proto.DeleteSessionRequest) (*proto.DeleteSessionResponse, error) {
	if err := database.DeleteSession(i.clients.Mongo, req); err != nil {
		return &proto.DeleteSessionResponse{
			Success: false,
		}, err
	}

	return &proto.DeleteSessionResponse{
		Success: true,
	}, nil
}

func (i *instance) FlushSessions(_ context.Context, req *proto.FlushSessionsRequest) (*proto.FlushSessionsResponse, error) {
	if err := database.FlushSessions(i.clients.Mongo, req); err != nil {
		return &proto.FlushSessionsResponse{
			Success: false,
		}, err
	}

	return &proto.FlushSessionsResponse{
		Success: true,
	}, nil
}

func (i *instance) ValidateEmail(_ context.Context, req *proto.ValidateEmailRequest) (*proto.ValidateEmailResponse, error) {
	if err := database.ValidateEmail(i.clients.Mongo, req); err != nil {
		return &proto.ValidateEmailResponse{
			Success: false,
		}, err
	}

	return &proto.ValidateEmailResponse{
		Success: true,
	}, nil
}

func (i *instance) CheckEmail(_ context.Context, req *proto.CheckEmailRequest) (*proto.CheckEmailResponse, error) {
	ok := database.CheckEmail(i.clients.Mongo, req)

	return &proto.CheckEmailResponse{
		Available: ok,
	}, nil
}
