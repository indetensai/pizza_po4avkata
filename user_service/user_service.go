package user_service

import (
	"context"
	"log"
	"net"
	"pizza/service"
	pb "pizza/service"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var client *mongo.Client

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {

}

func baza() {
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
}

func main() {
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
