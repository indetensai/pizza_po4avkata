package main

import (
	"log"
	"pizza/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var user_service service.UserServiceClient
var menu_service service.MenuServiceClient

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
	session_id, err := user_service.Login(c.Context(), &service.LoginRequest{
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"session_id": session_id.SessionId})
}

func get_menu_handler(c *fiber.Ctx) error {
	pizza, err := menu_service.GetMenu(c.Context(), &service.MenuRequest{})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"pizzas": pizza.Menu})
}

func main() {
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Post("user/login", login_handler)
	app.Post("menu/get", get_menu_handler)
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
	app.Listen(":8080")

}
