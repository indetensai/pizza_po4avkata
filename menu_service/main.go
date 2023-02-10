package main

import (
	"context"
	"log"
	"net"
	"os"
	"pizza/service"
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

type server struct {
	service.UnimplementedMenuServiceServer
}

type MenuPizza struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Ingredients []string           `bson:"ingredients"`
	Prices      map[string]int64   `bson:"prices"`
}

func (s *server) CreatePizza(
	ctx context.Context,
	request *service.PizzaCreatingRequest,
) (*service.PizzaCreatingResponse, error) {
	name := request.Content.GetName()
	description := request.Content.GetDescription()
	ingredients := request.Content.GetIngredients()
	prices := request.Content.GetPrices()
	_, err := db.Database("menu_service").Collection("pizza").InsertOne(
		ctx,
		bson.D{
			{Key: "name", Value: name},
			{Key: "description", Value: description},
			{Key: "ingredients", Value: ingredients},
			{Key: "prices", Value: prices}},
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &service.PizzaCreatingResponse{}, nil
}

func (s *server) GetMenu(
	ctx context.Context,
	request *service.MenuRequest,
) (*service.MenuResponse, error) {
	source, err := db.Database("menu_service").Collection("pizza").Find(
		ctx,
		bson.D{{}},
		nil,
	)
	if err != nil {
		return nil, err
	}
	var preresult []MenuPizza
	err = source.All(ctx, &preresult)
	if err != nil {
		return nil, err
	}
	var result []*service.MenuPizza
	for _, pizza := range preresult {
		result = append(result, &service.MenuPizza{
			Name:        pizza.Name,
			Description: pizza.Description,
			Ingredients: pizza.Ingredients,
			Prices:      pizza.Prices,
			Id:          pizza.ID.Hex()})
	}
	return &service.MenuResponse{Menu: result}, nil
}

func baza() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MENU_SERVICE_DATABASE")).
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
	ln, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}
	srv := grpc.NewServer()
	service.RegisterMenuServiceServer(srv, &server{})
	reflection.Register(srv)
	if err = srv.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
