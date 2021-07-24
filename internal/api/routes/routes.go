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

func UserDowngrade(route fiber.Router) {
	route.Post("", controllers.DowngradeMembership)
}

func ValidateKey(route fiber.Router) {
	route.Get("", controllers.AuthEthos)
}

func UserCheck(route fiber.Router) {
	route.Get("", controllers.Authenticate)
}
