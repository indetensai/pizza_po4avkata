package main

import (
	"log"
	"pizza/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var client service.UserServiceClient

func register_handler(c *fiber.Ctx) error {
	_, err := client.Register(c.Context(), &service.RegisterRequest{
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
	session_id, err := client.Login(c.Context(), &service.LoginRequest{
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"session_id": session_id.Session})
}

func main() {
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Post("user/login", login_handler)
	conn, err := grpc.Dial("localhost:443", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = service.NewUserServiceClient(conn)
	app.Listen(":8080")

}
