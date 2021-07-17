package routes

import (
	"ethos-dash/internal/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserCreate(route fiber.Router) {
	route.Post("", controllers.AddNewUser)
}

func UserUpdate(route fiber.Router) {
	route.Post("", controllers.UpdateMembership)
}

func UserUpgrade(route fiber.Router) {
	route.Post("", controllers.UpgradeSub)
}

func PartnershipRoutes(route fiber.Router) {
	route.Post("", controllers.NewPartner)
}
