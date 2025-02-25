package routes

import (
	"dgw-technical-test/handler"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, uh handler.UserHandler) {
	app.Post("/users/register", uh.Register)
	app.Post("/users/login", uh.Login)
}
