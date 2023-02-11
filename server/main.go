package main

import (
	"bytes"
	"encoding/json"
	"log"
	"pizza/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var user_service service.UserServiceClient
var menu_service service.MenuServiceClient
var order_service service.OrderServiceClient

type Order struct {
	Content []*service.OrderRequestPizza `bson:"content"`
	Session string                       `bson:"session"`
}

func register_handler(c *fiber.Ctx) error {
	_, err := user_service.Register(c.Context(), &service.RegisterRequest{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	})
	e, _ := status.FromError(err)
	if err != nil {
		if e.Code() == codes.AlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"note": "username already exists",
			})
		}
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func login_handler(c *fiber.Ctx) error {
	session, err := user_service.Login(c.Context(), &service.LoginRequest{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	})
	e, _ := status.FromError(err)
	if err != nil {
		if e.Code() == codes.Unauthenticated || e.Code() == codes.NotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"note": "invalid username or password",
			})
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"session": session.Session})
}

func menu_handler(c *fiber.Ctx) error {
	pizza, err := menu_service.GetMenu(c.Context(), &service.MenuRequest{})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"pizzas": pizza.Menu})
}

func order_handler(c *fiber.Ctx) error {
	var result Order
	err := json.NewDecoder(bytes.NewReader(c.Body())).Decode(&result)
	if err != nil {
		return err
	}
	order, err := order_service.Order(c.Context(), &service.OrderRequest{
		Content: result.Content,
		Session: result.Session,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"total bill": order.Bill})
}

func getting_orders_handler(c *fiber.Ctx) error {
	orders, err := order_service.GetOrders(c.Context(), &service.GetOrdersRequest{})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"orders": orders})
}

func main() {
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Post("/user/login", login_handler)
	app.Get("/menu", menu_handler)
	app.Post("/order", order_handler)
	app.Get("/orders", getting_orders_handler)
	conn, err := grpc.Dial("localhost:443", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	user_service = service.NewUserServiceClient(conn)
	connd, err := grpc.Dial("localhost:50001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	menu_service = service.NewMenuServiceClient(connd)
	connrd, err := grpc.Dial("localhost:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	order_service = service.NewOrderServiceClient(connrd)
	app.Listen(":8080")

}
