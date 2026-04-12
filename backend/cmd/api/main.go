package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/licandragon/FileTransfer/backend/internal/database"
	"github.com/licandragon/FileTransfer/backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Archivo .env no encontrado")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Fallo la conexion a la base de datos", err)
	}

	defer db.Close()

	app := fiber.New()
	routes.SetupRoutes(app, db)

	log.Fatal(app.Listen(":" + port))
}
