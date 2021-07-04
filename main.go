package main

import (
	"ethos-dash/internal/keygen"

	"github.com/gofiber/fiber/v2"

	"fmt"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	//app.Listen(":3000")

	fmt.Println(keygen.Keygen("abhaya#1149"))

}
