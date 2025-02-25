package routes

import (
	"dgw-technical-test/handler"
	"dgw-technical-test/middleware"

	_ "dgw-technical-test/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRoute(app *fiber.App, uh handler.UserHandler, bh handler.BookHandler) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	users := app.Group("/users")
	users.Post("/register", uh.Register)
	users.Post("/login", uh.Login)

	books := app.Group("/books", middleware.CustomJwtMiddleware())
	books.Post("/", bh.Create)
	books.Put("/:id", bh.Update)
	books.Delete("/:id", bh.Delete)
	books.Get("/", bh.FindAll)
	books.Get("/:id", bh.FindById)
}
