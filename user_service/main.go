package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"pizza/service"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var db *mongo.Client

type server struct {
	service.UnimplementedUserServiceServer
}

func (s *server) Register(
	ctx context.Context,
	request *service.RegisterRequest,
) (*service.RegisterResponse, error) {
	username := request.GetUsername()
	password := request.GetPassword()
	check := db.Database("user_service").Collection("users").FindOne(
		ctx,
		bson.D{{Key: "username", Value: username}},
		nil,
	)
	if check != nil {
		return nil, status.Error(codes.AlreadyExists, "username already taken")
	}
	_, err := db.Database("user_service").Collection("users").InsertOne(
		ctx,
		bson.D{{Key: "username", Value: username}, {Key: "password", Value: password}},
		nil,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &service.RegisterResponse{}, nil

}

func baza() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("USER_SERVICE_DATABASE")).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	db, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	godotenv.Load()
	a := (os.Getenv("USER_SERVICE_DATABASE"))
	fmt.Println(a)
	baza()
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	service.RegisterUserServiceServer(srv, &server{})
	reflection.Register(srv)
	if err = srv.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
