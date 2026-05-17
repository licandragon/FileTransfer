package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/licandragon/FileTransfer/backend/internal/models"
	"github.com/licandragon/FileTransfer/backend/internal/services"
)

// ------------------------------------------------------------------
// DTOs de entrada
// ------------------------------------------------------------------

type CreateTransferRequest struct {
	SenderEmail  string     `json:"sender_email"`
	SubjectEmail string     `json:"subject_email"`
	MessageEmail string     `json:"message_email"`
	Recipients   []string   `json:"recipients"`
	TotalFiles   int        `json:"total_files"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

// ------------------------------------------------------------------
// DTOs de respuesta
// ------------------------------------------------------------------

type CreateTransferResponse struct {
	UploadToken    uuid.UUID `json:"upload_token"`
	StatusTransfer string    `json:"status_transfer"`
}

type TransferStatusResponse struct {
	CompletedIndices []int `json:"completedIndices"`
}

type FileUploadResponse struct {
	FileIndex    int    `json:"file_index"`
	OriginalName string `json:"original_name"`
	StatusFile   string `json:"status_file"`
}

type CompleteTransferResponse struct {
	DownloadToken uuid.UUID `json:"download_token"`
}

type TransferDetailResponse struct {
	ID             string       `json:"id,omitempty"`
	DownloadToken  string       `json:"download_token"`
	SenderEmail    string       `json:"sender_email"`
	SubjectEmail   string       `json:"subject_email"`
	MessageEmail   string       `json:"message_email"`
	Recipients     []string     `json:"recipients"`
	ExpiresAt      *time.Time   `json:"expires_at,omitempty"`
	StatusTransfer string       `json:"status_transfer"`
	TotalFiles     int          `json:"file_count,omitempty"`
	Files          []FileDetail `json:"files"`
	CreatedAt      time.Time    `json:"created_at"`
}

type TransferDownloadResponse struct {
	SenderEmail  string       `json:"sender_email"`
	SubjectEmail string       `json:"subject_email"`
	MessageEmail string       `json:"message_email"`
	ExpiresAt    *time.Time   `json:"expires_at,omitempty"`
	Files        []FileDetail `json:"files"`
}

type FileDetail struct {
	FileIndex    int    `json:"file_index"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
}

// ------------------------------------------------------------------
// Handler
// ------------------------------------------------------------------

type TransferHandler struct {
	service services.TransferService
}

func NewTransferHandler(service services.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

// POST /api/transfers
func (h *TransferHandler) CreateTransfer(c fiber.Ctx) error {
	body := c.Body()
	log.Printf("Cuerpo recibido: %s", string(body))
	log.Printf("CT: %s, Body: %q", c.Get("Content-Type"), string(c.Body()))

	var req CreateTransferRequest
	log.Printf("Cuerpo: %s", body)
	if err := c.Bind().JSON(&req); err != nil {
		log.Printf("Error binding: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}
	fmt.Println(req.SenderEmail)
	fmt.Println(req.TotalFiles)
	if req.SenderEmail == "" || req.TotalFiles <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Faltan campos obligatorios"})
	}

	if req.ExpiresAt == nil {
		exp := time.Now().Add(7 * 24 * time.Hour)
		req.ExpiresAt = &exp
	}

	transfer := &models.Transfer{
		SenderEmail:  req.SenderEmail,
		SubjectEmail: req.SubjectEmail,
		MessageEmail: req.MessageEmail,
		Recipients:   req.Recipients,
		TotalFiles:   req.TotalFiles,
		ExpiresAt:    req.ExpiresAt,
	}

	result, err := h.service.CreateTransfer(c.Context(), transfer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(CreateTransferResponse{
		UploadToken:    result.UploadToken,
		StatusTransfer: result.StatusTransfer,
	})
}

// POST /api/transfer/:uploadToken/files
func (h *TransferHandler) AddFile(c fiber.Ctx) error {
	fmt.Println(c)
	uploadTokenStr := c.Params("uploadToken")
	uploadToken, err := uuid.Parse(uploadTokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "upload_token inválido"})
	}

	fileIndexStr := c.FormValue("file_index")
	if fileIndexStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "file_index es requerido"})
	}
	fileIndex, err := strconv.Atoi(fileIndexStr)
	if err != nil || fileIndex < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "file_index inválido"})
	}
	//Test para probar error en subida de algun archivo
	/*if fileIndex == 1 {
		return c.Status(500).JSON(fiber.Map{"error": "Simulación: Conexión perdida con el almacenamiento"})
	}*/

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Archivo no encontrado"})
	}

	uploaded, err := h.service.AddFileToTransfer(c.Context(), fileHeader, uploadToken, fileIndex)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(FileUploadResponse{
		FileIndex:    uploaded.FileIndex,
		OriginalName: uploaded.OriginalName,
		StatusFile:   uploaded.StatusFile,
	})
}

// GET /api/transfer/:uploadToken/status
func (h *TransferHandler) GetUploadStatus(c fiber.Ctx) error {
	uploadTokenStr := c.Params("uploadToken")
	uploadToken, err := uuid.Parse(uploadTokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "upload_token inválido"})
	}

	indices, err := h.service.GetTransferUploadStatus(c.Context(), uploadToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(TransferStatusResponse{
		CompletedIndices: indices,
	})
}

// PATCH /api/transfer/:uploadToken/complete
func (h *TransferHandler) CompleteTransfer(c fiber.Ctx) error {
	uploadTokenStr := c.Params("uploadToken")
	uploadToken, err := uuid.Parse(uploadTokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "upload_token inválido"})
	}

	downloadToken, err := h.service.CompleteTransfer(c.Context(), uploadToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(CompleteTransferResponse{
		DownloadToken: downloadToken,
	})
}

// GET /api/download/:downloadToken
func (h *TransferHandler) DownloadInfo(c fiber.Ctx) error {
	downloadTokenStr := c.Params("downloadToken")
	downloadToken, err := uuid.Parse(downloadTokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "upload_token inválido"})
	}

	transfer, err := h.service.GetTransferByDownloadToken(c.Context(), downloadToken)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transferencia no encontrada"})
	}

	files := make([]FileDetail, len(transfer.Files))
	for i, f := range transfer.Files {
		files[i] = FileDetail{
			FileIndex:    f.FileIndex,
			OriginalName: f.OriginalName,
			Size:         f.SizeFile,
			MimeType:     f.MimeType,
		}
	}

	return c.Status(200).JSON(TransferDownloadResponse{
		SenderEmail:  transfer.SenderEmail,
		SubjectEmail: transfer.SubjectEmail,
		MessageEmail: transfer.MessageEmail,
		ExpiresAt:    transfer.ExpiresAt,
		Files:        files,
	})
}

// GET /api/download/:downloadToken/files/:fileIndex
func (h *TransferHandler) DownloadFile(c fiber.Ctx) error {
	downloadTokenStr := c.Params("downloadToken")
	downloadToken, err := uuid.Parse(downloadTokenStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "upload_token inválido"})
	}

	fileIndexStr := c.Params("fileIndex")
	fileIndex, err := strconv.Atoi(fileIndexStr)
	if err != nil || fileIndex < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "índice de archivo inválido"})
	}

	transfer, err := h.service.GetTransferByDownloadToken(c.Context(), downloadToken)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transferencia no encontrada"})
	}

	// Buscar el archivo por índice
	var targetFile *models.File
	for i := range transfer.Files {
		if transfer.Files[i].FileIndex == fileIndex {
			targetFile = &transfer.Files[i]
			break
		}
	}
	if targetFile == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Archivo no encontrado"})
	}

	// Obtener URL firmada
	signedURL, err := h.service.GetFileSignedURL(c.Context(), downloadToken, fileIndex)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Descargar el archivo desde Supabase
	resp, err := http.Get(signedURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener archivo"})
	}
	defer resp.Body.Close()

	// Configurar headers para forzar descarga
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, targetFile.OriginalName))
	c.Set("Content-Type", resp.Header.Get("Content-Type"))
	if resp.ContentLength > 0 {
		c.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}

	// Transmitir el archivo al cliente sin cargarlo en memoria
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al transmitir archivo"})
	}

	return nil
}

// ------------------------------------------------------------------
// Placeholders para futuros endpoints con autenticación
// ------------------------------------------------------------------

func (h *TransferHandler) ListUserTransfers(c fiber.Ctx) error {
	return c.Status(501).JSON(fiber.Map{"error": "no implementado"})
}

func (h *TransferHandler) GetTransfer(c fiber.Ctx) error {
	return c.Status(501).JSON(fiber.Map{"error": "no implementado"})
}

func (h *TransferHandler) DeleteTransfer(c fiber.Ctx) error {
	return c.Status(501).JSON(fiber.Map{"error": "no implementado"})
}
