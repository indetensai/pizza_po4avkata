package main

import (
	"context"
	"log"
	"net"
	"os"
	"pizza/service"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db *mongo.Client

var user_service service.UserServiceClient

type server struct {
	service.UnimplementedOrderServiceServer
}

func (s *server) Order(
	ctx context.Context,
	request *service.OrderRequest,
) (*service.OrderResponse, error) {
	pizzas := request.GetContent()
	session := request.GetSession()
	username, err := user_service.SessionCheck(
		ctx,
		&service.SessionCheckRequest{
			Session: session,
		},
	)
	if err != nil || username.Username == "" {
		return nil, err
	}
	var bill int64
	var result service.Pizza
	for _, pizza := range pizzas {
		object_id, err := primitive.ObjectIDFromHex(pizza.PizzaId)
		if err != nil {
			return nil, err
		}
		err = db.Database("menu_service").Collection("pizza").FindOne(
			ctx,
			bson.D{{Key: "_id", Value: object_id}},
		).Decode(&result)
		if err != nil {
			return nil, err
		}
		bill += result.Prices[pizza.Size]
	}
	_, err = db.Database("order_service").Collection("orders").InsertOne(
		ctx,
		bson.D{
			{Key: "content", Value: pizzas},
			{Key: "status", Value: "awaiting"},
			{Key: "time", Value: time.Now().String()},
			{Key: "username", Value: username.Username},
			{Key: "bill", Value: bill},
		},
	)
	if err != nil {
		return nil, err
	}
	return &service.OrderResponse{Bill: strconv.Itoa(int(bill))}, nil
}

func baza() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("ORDER_SERVICE_DATABASE")).
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
	baza()
	ln, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	service.RegisterOrderServiceServer(srv, &server{})
	reflection.Register(srv)
	conn, err := grpc.Dial("localhost:443", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	user_service = service.NewUserServiceClient(conn)
	if err = srv.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
