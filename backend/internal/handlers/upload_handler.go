package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/licandragon/FileTransfer/backend/internal/models"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
)

type UploadHandler struct {
	TransferRepo *repository.TransferRepository
	FileRepo     *repository.FileRepository
}

func NewUploadHandler(tr *repository.TransferRepository, fr *repository.FileRepository) *UploadHandler {
	return &UploadHandler{
		TransferRepo: tr,
		FileRepo:     fr,
	}
}

func (h *UploadHandler) Upload(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Error al leer el formulario",
		})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "No se enviaron archivos",
		})
	}

	ctx := context.Background()
	// 1. Crear Transfer
	transfer := models.Transfer{
		ID:            uuid.New().String(),
		DownloadToken: uuid.New().String(),
		CreatedAt:     time.Now(),
	}

	err = h.TransferRepo.Create(ctx, &transfer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create transfer",
		})
	}

	return c.JSON(fiber.Map{
		"status":         "complete",
		"download_token": transfer.DownloadToken,
	})

}
