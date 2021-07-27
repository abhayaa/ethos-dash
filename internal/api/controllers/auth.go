package controllers

import (
	"ethos-dash/internal/db"
	"ethos-dash/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func DiscordAuth(c *fiber.Ctx) error {
	c.Accepts("application/json")

	type Request struct {
		Code string `json:"code"`
		Key  string `json:"key"`
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

	var user *utils.User = utils.GetDiscordInfo(body.Code)
	var found bool = db.UserCheck(user.Id)

	if found {
		q := db.GetUserById(user.Id)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":    true,
			"found":      found,
			"username":   q.Username,
			"avatar":     user.Avatar,
			"expiration": q.Expiration,
			"ethosKey":   q.EthosKey,
			"plan":       q.PlanType,
			"memberType": q.MemberType,
			"userId":     q.UserId,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"username": user.Username + "#" + user.Discriminator,
		"avatar":   user.Avatar,
		"found":    found,
	})
}
