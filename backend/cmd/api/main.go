package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/licandragon/FileTransfer/backend/internal/routes"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
