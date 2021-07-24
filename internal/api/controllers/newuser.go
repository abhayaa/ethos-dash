package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/keygen"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func AddNewUser(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		UserId string `json:"userId"`
		Key    string `json:"key"`
		Email  string `json:"email"`
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
			"apiKey":  utils.GetEnvKey("API_KEY"),
			"key":     body.Key,
		})
	}

	user := db.User{
		UserId:   body.UserId,
		EthosKey: keygen.Keygen(body.UserId),
		Email:    body.Email,
	}

	error := db.InsertUser(user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": error,
	})

}
