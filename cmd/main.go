package main

import (
	"dgw-technical-test/config"
	"dgw-technical-test/handler"
	"dgw-technical-test/repository"
	"dgw-technical-test/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

// @title           DGW-Technical-Test
// @version         1.0
// @description     This is a simple crud api
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /

func main() {
	app := fiber.New()
	app.Use(logger.New())

	db := config.NewDatabase()
	validate := validator.New()
	userRepository := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepository, validate)

	bookRepository := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepository, validate)

	routes.NewRoute(app, *userHandler, *bookHandler)

	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	port := os.Getenv("PORT")
	go func() {
		if err := app.Listen(":" + port); err != nil {
			errChan <- err
		}
	}()

	defer func() {
		log.Println("Closing database connection...")
		db.Close()

		log.Println("Shutting down Fiber server...")
		app.Shutdown()
	}()

	select {
	case err := <-errChan:
		log.Printf("Fail to serve: %v\n", err)
	case <-stopChan:
		log.Println("Received shutdown signal...")
	}
}
