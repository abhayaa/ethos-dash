package main

import (
	"ethos-dash/internal/api/routes"
	_ "ethos-dash/internal/keygen"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	setupRoutes(app)

	err := app.Listen(":8000")

	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":8000"))

}

func setupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "API is up/running",
		})
	})

	userapi := app.Group("/users")
	partnerapi := app.Group("/partner")
	discordapi := app.Group("/discord")
	adminapi := app.Group("/admin")

	// user routes
	routes.UserCreate(userapi.Group("/createUser"))
	routes.UserUpdate(userapi.Group("/update"))
	routes.UserUpgrade(userapi.Group("/upgrade"))
	routes.UserDowngrade(userapi.Group("/downgrade"))
	routes.UserCheck(userapi.Group("/usercheck"))

	// admin routes
	routes.GenerateKey(adminapi.Group("/generatekey"))
	//partner routes
	routes.PartnershipRoutes(partnerapi.Group("/newpartner"))
	routes.ValidateKey(partnerapi.Group("/auth"))

	//discord routes
	routes.DiscordAuth(discordapi.Group("/auth"))
}
