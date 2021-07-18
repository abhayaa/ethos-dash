package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func DowngradeMembership(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		UserId string `json:"userId"`
		Bot    string `json:"bot"`
		Key    string `json:"key"`
		BotKey string `json:"botKey"`
	}

	var body Request
	err := c.BodyParser(&body)

	if err != nil {
		log.Printf("error while parsing JSON")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "cannot parse JSON",
		})
	}

	if body.Key != utils.GetEnvKey("API_KEY") {
		log.Printf("API key does not match")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
		})
	}

	botDowngrade := db.BotKey{
		BotKey: body.BotKey,
		Bot:    body.Bot,
		UserId: body.UserId,
	}

	error := db.DowngradeKey(botDowngrade)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": error,
	})
}
