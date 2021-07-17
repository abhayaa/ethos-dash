package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func NewPartner(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		BotName string `json:"BotKey"`
		Key     string `json:"key"`
	}

	var body Request
	err := c.BodyParser(&body)

	if err != nil {
		log.Printf("Error while parsing json")
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

	error := db.CreateBotTable(body.BotName)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": error,
	})
}
