package routes

import (
	"dgw-technical-test/handler"
	"dgw-technical-test/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, uh handler.UserHandler, bh handler.BookHandler) {
	users := app.Group("/users")
	users.Post("/register", uh.Register)
	users.Post("/login", uh.Login)

	books := app.Group("/books", middleware.CustomJwtMiddleware())
	books.Post("/", bh.Create)
	books.Put("/", bh.Update)
	books.Delete("/", bh.Delete)
	books.Get("/", bh.FindAll)
	books.Get("/:id", bh.FindById)
}
