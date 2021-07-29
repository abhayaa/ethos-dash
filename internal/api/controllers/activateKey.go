package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ActivateKey(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
		EthosKey     string `json:"ethosKey"`
		Key          string `json:"key"`
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

	discordAuthRes := utils.GetDiscordInfo(body.AccessToken)

	var user db.User
	user.AccessToken = body.AccessToken
	user.RefreshToken = body.RefreshToken
	user.EthosKey = body.EthosKey

	user.UserId = discordAuthRes.Id
	user.Email = discordAuthRes.Email
	user.Username = discordAuthRes.Username + "#" + discordAuthRes.Discriminator

	error := db.InsertUser(user)

	keyValidate := db.ValidateKetActivation(body.EthosKey)

	if !keyValidate {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "invalid key",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": error,
	})
}
