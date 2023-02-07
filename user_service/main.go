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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var client *mongo.Client

type server struct {
	service.UnimplementedUserServiceServer
}

//func (s *server) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {

//}

func baza() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("USER_SERVICE_DATABASE")).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
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
