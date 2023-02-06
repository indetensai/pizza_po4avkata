package user_service

import (
	"context"
	"log"
	pb "pizza/service"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://indetensai:<password>@cluster0.wzkqiif.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return &pb.RegisterResponse{}, nil
}
