package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func AuthEthos(c *fiber.Ctx) error {
	c.Accepts("json", "text")

	type Request struct {
		BotName string `json:"botName"`
		Key     string `json:"key"`
		BotKey  string `json:"botKey"`
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

	exp, exist := db.ValidateKey(body.BotKey, body.BotName)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"exists": exist,
		"exp":    exp,
	})
}

/**
{
    "BotName": "stellar",
    "Key" : "nplq~VL}W3[3'p2']gpmZF*U=V+^-!",
    "BotKey": "SLAYWC-RUIV-M5SPNBS-5SPN"
}
*/
