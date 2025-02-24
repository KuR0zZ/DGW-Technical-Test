package main

import (
	"dgw-technical-test/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	_ = config.NewDatabase()
}
