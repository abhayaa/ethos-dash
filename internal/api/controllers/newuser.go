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
		UserId       string `json:"userId"`
		Key          string `json:"key"`
		Email        string `json:"email"`
		AccessToken  string `json:"accesstoken"`
		RefreshToken string `json:"refreshtoken"`
		Username     string `json:"username"`
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
			"key":     body.Key,
		})
	}

	user := db.User{
		UserId:       body.UserId,
		EthosKey:     keygen.Keygen(body.UserId),
		Email:        body.Email,
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		Username:     body.Username,
	}

	error := db.InsertUser(user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": error,
	})

}
