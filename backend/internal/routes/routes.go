package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/handlers"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
)

// Aqui se definiran los endpoints de la API
func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {

	transferRepo := repository.NewTransferRepository(db)
	fileRepo := repository.NewFileRepository(db)

	uploadHandler := handlers.NewUploadHandler(transferRepo, fileRepo)

	app.Post("/upload", uploadHandler.Upload)

	app.Post("/download/:token", func(c fiber.Ctx) error {
		return c.SendString("Descargando archivo")
	})

	app.Get("/file/:token", func(c fiber.Ctx) error {
		return c.SendString("Obteninedo archivo")
	})

	app.Delete("/file/:token", func(c fiber.Ctx) error {
		return c.SendString("Eliminando archivo")
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("API Running")
	})
}
