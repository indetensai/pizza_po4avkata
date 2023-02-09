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

var client service.UserServiceClient

func main() {
	app := fiber.New()
	app.Post("/user/register", register_handler)
	conn, err := grpc.Dial("localhost:443", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client = service.NewUserServiceClient(conn)
	app.Listen(":8080")

}
