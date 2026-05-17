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
	handler := handlers.NewTransferHandler(transferService)

	api := app.Group("/api")

	//Grupo para la Creacion de la transferencia y subida de archivos a la misma
	upload := api.Group("/transfer")

	//Se crea transferencia (metadatos + file_count) -> devuelve upload_token
	upload.Post("/init", handler.CreateTransfer)

	// Operaciones con upload_token (subida y completitud)

	upload.Post("/:uploadToken/file", handler.AddFile)

	api.Get("/transfers/:uploadToken/status", handler.GetUploadStatus)

	// Se finaliza la transferencia  y retorna
	upload.Patch(":uploadToken/complete", handler.CompleteTransfer)

	//Endpoind para url publica para
	api.Get("/download/:downloadToken", handler.DownloadInfo)

	// Endpoint para descargar un file en especifico de la transferencia
	api.Post("/download/:downloadToken/files/:fileIndex", handler.DownloadFile)

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("API Running")
	})
}
