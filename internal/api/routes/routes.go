package routes

import (
	"ethos-dash/internal/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserCreateRoute(route fiber.Router) {
	route.Post("", controllers.AddNewUser)
}
