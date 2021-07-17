package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func UpdateMembership(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		UserId string `json:"userId"`
		Key    string `json:"key"`
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

	error := db.UpdateMembership(body.UserId)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": error,
	})
}
