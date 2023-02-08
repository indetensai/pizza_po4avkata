package main

import (
	"log"
	"pizza/service"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func register_handler(c *fiber.Ctx) error {
	_, err := client.Register(c.Context(), &service.RegisterRequest{Username: c.FormValue("username"), Password: c.FormValue("password")})
	return err
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
