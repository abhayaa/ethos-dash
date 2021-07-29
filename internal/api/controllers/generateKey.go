package controllers

import (
	"log"
	"math/rand"

	"ethos-dash/internal/db"
	"ethos-dash/internal/keygen"
	"ethos-dash/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GenerateKey(c *fiber.Ctx) error {

	c.Accepts("json", "text")

	type Request struct {
		Key       string `json:"key"`
		Generator string `json:"generator"`
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

	if !utils.ValidateApiKey(body.Key) {
		log.Printf("API key does not match")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
		})
	}

	key := keygen.Keygen(utils.GenerateRandomString(rand.Intn(100)))

	insertKey := db.AddKey(key, body.Generator)
	if !insertKey {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"key":     key,
	})
}
