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
		Id  string `json:"Id"`
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

	if !utils.ValidateApiKey(body.Key) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"key":     body.Key,
		})
	}

	status := db.UserCheck(body.Id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"found":   status,
	})
}
