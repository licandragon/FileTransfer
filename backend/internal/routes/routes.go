package routes

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/handlers"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
	"github.com/licandragon/FileTransfer/backend/internal/services"
	"github.com/licandragon/FileTransfer/backend/internal/storage"
)

// Aqui se definiran los endpoints de la API
func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {

	url := os.Getenv("SUPABASE_URL")
	api_key := os.Getenv("SUPABASE_SERVICE_KEY")

	transfer := repository.NewTransferRepository(db)
	storage := storage.NewSupabaseStorage(url, api_key)

	transferService := services.NewTransferService(transfer, storage)

	handler := handlers.NewUploadHandler(transferService)

	app.Post("/upload", handler.Upload)

	app.Get("/download/:token", func(c fiber.Ctx) error {
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
