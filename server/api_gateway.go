package server

import (
	"github.com/gofiber/fiber/v2"
)

func register_handler(c *fiber.Ctx) error {

}

func main() {
	app := fiber.New()
	app.Post("/user/register", register_handler)
	app.Listen(":8080")
}
