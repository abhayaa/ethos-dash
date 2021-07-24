package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		Id  string `json:"userId"`
		Key string `json:"key"`
	}

	var body Request
	err := c.BodyParser(&body)

	if err != nil {
		log.Printf("Error parsing json")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "cannot parse JSON",
		})
	}

	if body.Key != utils.GetEnvKey("API_KEY") {
		log.Printf("API key does not match")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"apiKey":  utils.GetEnvKey("API_KEY"),
			"key":     body.Key,
		})
	}

	status := db.UserCheck(body.Id)

	return c.Status(fiber.StatusFound).JSON(fiber.Map{
		"success": true,
		"found":   status,
	})
}
