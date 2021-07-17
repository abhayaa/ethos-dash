package main

import (
	"ethos-dash/internal/api/routes"
	_ "ethos-dash/internal/keygen"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)

	err := app.Listen(":8000")

	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":3000"))

}

func setupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "At the endpoint",
		})
	})

	userapi := app.Group("/users")
	partnerapi := app.Group("/partner")

	userapi.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint",
		})
	})

	routes.UserCreate(userapi.Group("/createUser"))
	routes.UserUpdate(userapi.Group("/update"))
	routes.UserUpgrade(userapi.Group("/upgrade"))

	routes.PartnershipRoutes(partnerapi.Group("/newpartner"))
}
