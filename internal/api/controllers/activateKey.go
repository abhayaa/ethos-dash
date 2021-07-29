package controllers

// import (
// 	"ethos-dash/internal/db"
// 	"ethos-dash/internal/utils"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// )

// func ActivateKey(c *fiber.Ctx) error {
// 	c.Accepts("json", "text")

// 	type Request struct {
// 		Code string `json:"code"`
// 		EthosKey string `json:"ethosKey"`
// 		Key string `json:"key"`
// 	}

// 	var body Request
// 	err := c.BodyParser(&body)

// 	if err != nil {
// 		log.Printf("Error while parsing json")
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"success": false,
// 			"message": "cannot parse JSON",
// 		})
// 	}

// 	if !utils.ValidateApiKey(body.Key) {
// 		log.Printf("API key does not match")
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"success": false,
// 			"key":     body.Key,
// 		})
// 	}

// 	discordAuthRes := utils.GetDiscordInfo(body.Code)

// 	var user db.User
// 	user.UserId = discordAuthRes.Id
// 	user.EthosKey = body.EthosKey
// 	user.Email = discordAuthRes.Email
// 	user.AccessToken = discordAuthRes.

// 	error := db.InsertUser(user)

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"success": true,
// 		"message": error,
// 	})
// }
