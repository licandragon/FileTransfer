package services

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/licandragon/FileTransfer/backend/internal/models"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
	"github.com/licandragon/FileTransfer/backend/internal/storage"
)

// Se define interfaz del servicio
type TransferService interface {
	CreateTransfer(ctx context.Context, transfer *models.Transfer, files []*multipart.FileHeader) (*models.Transfer, error)
	GetTransferByToken(ctx context.Context, token string) (*models.Transfer, error)
}

type transferService struct {
	repo    repository.TransferRepository
	storage storage.FileStorage
}

func NewTransferService(repo repository.TransferRepository, storage storage.FileStorage) TransferService {
	return &transferService{
		repo:    repo,
		storage: storage,
	}
}

// CreateTransfer implements [TransferService].
func (t *transferService) CreateTransfer(ctx context.Context, transfer *models.Transfer, files []*multipart.FileHeader) (*models.Transfer, error) {
	transfer.DownloadToken = uuid.New().String()
	bucketName := "transfers"

	var uploadedFiles []models.File

	// 2. Subir archivos a la nube
	for _, fileHeader := range files {
		// Creamos un nombre de archivo único para evitar colisiones en el storage
		storagePath := fmt.Sprintf("%s/%s", transfer.DownloadToken, fileHeader.Filename)

		_, err := t.storage.UploadFile(ctx, bucketName, storagePath, fileHeader)
		if err != nil {
			return nil, fmt.Errorf("error subiendo archivo %s: %w", fileHeader.Filename, err)
		}

		// Preparar el modelo de File para la DB
		newFile := models.File{
			Filename:     fileHeader.Filename,
			OriginalName: fileHeader.Filename,
			Size:         fileHeader.Size,
			MimeType:     fileHeader.Header.Get("Content-Type"),
			StoragePath:  storagePath,
			Bucket:       bucketName,
		}
		uploadedFiles = append(uploadedFiles, newFile)
	}

	transfer.Files = uploadedFiles

	// Guardado transaccional en DB usando el Repo
	if err := t.repo.Create(ctx, transfer); err != nil {
		return nil, fmt.Errorf("falló persistencia en DB: %w", err)
	}

	return transfer, nil
}

// GetTransferByToken implements [TransferService].
func (t *transferService) GetTransferByToken(ctx context.Context, token string) (*models.Transfer, error) {
	return t.repo.GetByToken(ctx, token)
}
