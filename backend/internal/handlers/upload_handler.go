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

	senderEmail := ""
	if val, ok := form.Value["sender_email"]; ok && len(val) > 0 {
		senderEmail = val[0]
	}

	subjectEmail := ""
	if val, ok := form.Value["subject_email"]; ok && len(val) > 0 {
		subjectEmail = val[0]
	}

	var messagePtr string
	if val, ok := form.Value["message_email"]; ok && len(val) > 0 && val[0] != "" {
		msg := val[0]
		messagePtr = msg
	}

	recipients := []string{}
	if val, ok := form.Value["recipients"]; ok {
		recipients = val
	}

	ctx := context.Background()
	expiration := time.Now().Add(7 * 24 * time.Hour)
	// 1. Crear Transfer
	transfer := models.Transfer{
		DownloadToken: uuid.New().String(),
		SenderEmail:   senderEmail,
		SubjectEmail:  subjectEmail,
		MessageEmail:  messagePtr,
		Recipients:    recipients,
		UserID:        nil,         // O un ID de usuario por defecto si es obligatorio
		ExpiresAt:     &expiration, // Expira en 7 días
	}

	err = h.TransferRepo.Create(ctx, &transfer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error DB: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":         "complete",
		"download_token": transfer.DownloadToken,
	})

}
