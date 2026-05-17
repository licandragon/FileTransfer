package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/licandragon/FileTransfer/backend/internal/models"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
	"github.com/licandragon/FileTransfer/backend/internal/storage"
)

// Se define interfaz del servicio
type TransferService interface {
	CreateTransfer(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	AddFileToTransfer(ctx context.Context, fileHeader *multipart.FileHeader, uploadToken uuid.UUID, fileIndex int) (*models.File, error)
	GetTransferUploadStatus(ctx context.Context, uploadToken uuid.UUID) ([]int, error)
	CompleteTransfer(ctx context.Context, uploadToken uuid.UUID) (uuid.UUID, error)
	GetTransferByUploadToken(ctx context.Context, uploadToken uuid.UUID) (*models.Transfer, error)
	GetTransferByDownloadToken(ctx context.Context, downloadToken uuid.UUID) (*models.Transfer, error)
	GetTransferByID(ctx context.Context, id uuid.UUID) (*models.Transfer, error)
	ListUserTransfers(ctx context.Context, userID uuid.UUID) ([]models.Transfer, error)
	DeleteTransfer(ctx context.Context, transferID, userID uuid.UUID) error
	GetFileSignedURL(ctx context.Context, downloadToken uuid.UUID, fileIndex int) (string, error)
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

// CreateTransfer crea una transferencia vacía y devuelve el upload_token (y otros datos).
func (s *transferService) CreateTransfer(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error) {
	if transfer.SenderEmail == "" {
		return nil, fmt.Errorf("el email del remitente es obligatorio")
	}
	if transfer.TotalFiles <= 0 {
		return nil, fmt.Errorf("file_count debe ser mayor que 0")
	}

	if err := s.repo.CreateTransfer(ctx, transfer); err != nil {
		return nil, fmt.Errorf("error al crear transferencia: %w", err)
	}
	return transfer, nil
}

// AddFileToTransfer sube un archivo a Supabase y lo registra en la BD usando el upload_token.
func (s *transferService) AddFileToTransfer(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	uploadToken uuid.UUID,
	fileIndex int,
) (*models.File, error) {
	transfer, err := s.repo.GetByUploadToken(ctx, uploadToken)
	if err != nil {
		return nil, fmt.Errorf("transferencia no encontrada: %w", err)
	}

	if transfer.StatusTransfer != "pending" {
		return nil, fmt.Errorf("la transferencia no acepta más archivos (estado: %s)", transfer.StatusTransfer)
	}
	fmt.Print("fileIndex", fileIndex, "y transferTotal: ", transfer.TotalFiles)
	if fileIndex < 0 || fileIndex >= transfer.TotalFiles {
		return nil, fmt.Errorf("índice de archivo %d fuera del rango esperado [0-%d]", fileIndex, transfer.TotalFiles-1)
	}

	safeFilename := sanitizeFilename(fileHeader.Filename)
	storagePath := fmt.Sprintf("%s/%s_%s", transfer.ID, uuid.New().String(), safeFilename)
	bucketName := "transfers"

	_, err = s.storage.UploadFile(ctx, bucketName, storagePath, fileHeader)
	if err != nil {
		return nil, fmt.Errorf("error subiendo archivo: %w", err)
	}

	file := &models.File{
		TransferID:   transfer.ID,
		FileIndex:    fileIndex,
		Filename:     safeFilename,
		OriginalName: fileHeader.Filename,
		SizeFile:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		StoragePath:  storagePath,
		Bucket:       bucketName,
	}

	_, err = s.repo.AddOrRetryFile(ctx, file)
	if err != nil {
		_ = s.storage.DeleteFile(ctx, bucketName, storagePath)
		return nil, fmt.Errorf("error al registrar archivo: %w", err)
	}

	return file, nil
}

func (s *transferService) GetTransferUploadStatus(ctx context.Context, uploadToken uuid.UUID) ([]int, error) {
	// Se valida que la transferencia exista
	_, err := s.repo.GetByUploadToken(ctx, uploadToken)
	if err != nil {
		return nil, fmt.Errorf("transferencia no encontrada: %w", err)
	}

	// Se obtienen los indices completados
	indices, err := s.repo.GetUploadedIndices(ctx, uploadToken)
	if err != nil {
		return nil, fmt.Errorf("error al obtener índices subidos: %w", err)
	}

	return indices, nil
}

// CompleteTransfer finaliza la transferencia y devuelve el download_token.
func (s *transferService) CompleteTransfer(ctx context.Context, uploadToken uuid.UUID) (uuid.UUID, error) {
	transfer, err := s.repo.GetByUploadToken(ctx, uploadToken)
	if err != nil {
		return uuid.Nil, err
	}
	downloadToken, err := s.repo.CompleteTransfer(ctx, transfer.ID, uploadToken)
	if err != nil {
		return uuid.Nil, err
	}
	return downloadToken, nil
}

// GetTransferByUploadToken obtiene la transferencia (incluye archivos) para gestión.
func (s *transferService) GetTransferByUploadToken(ctx context.Context, uploadToken uuid.UUID) (*models.Transfer, error) {
	return s.repo.GetByUploadToken(ctx, uploadToken)
}

// GetTransferByDownloadToken obtiene la transferencia pública para descarga (solo si está 'complete').
func (s *transferService) GetTransferByDownloadToken(ctx context.Context, downloadToken uuid.UUID) (*models.Transfer, error) {
	return s.repo.GetByDownloadToken(ctx, downloadToken)
}

// GetTransferByID obtiene cualquier transferencia por su ID interno.
func (s *transferService) GetTransferByID(ctx context.Context, id uuid.UUID) (*models.Transfer, error) {
	return s.repo.GetByID(ctx, id)
}

// ListUserTransfers lista las transferencias de un usuario autenticado.
func (s *transferService) ListUserTransfers(ctx context.Context, userID uuid.UUID) ([]models.Transfer, error) {
	return s.repo.ListByUser(ctx, userID)
}

// DeleteTransfer elimina una transferencia pendiente y sus archivos asociados.
func (s *transferService) DeleteTransfer(ctx context.Context, transferID, userID uuid.UUID) error {
	transfer, err := s.repo.GetByID(ctx, transferID)
	if err != nil {
		return err
	}
	if transfer.UserID == nil || *transfer.UserID != userID {
		return fmt.Errorf("no autorizado para eliminar esta transferencia")
	}
	for _, f := range transfer.Files {
		_ = s.storage.DeleteFile(ctx, f.Bucket, f.StoragePath)
	}
	return s.repo.DeletePendingTransfer(ctx, transferID)
}

// GetFileSignedURL genera una URL firmada para descargar un archivo concreto.
func (s *transferService) GetFileSignedURL(ctx context.Context, downloadToken uuid.UUID, fileIndex int) (string, error) {
	transfer, err := s.repo.GetByDownloadToken(ctx, downloadToken)
	if err != nil {
		return "", fmt.Errorf("transferencia no encontrada: %w", err)
	}

	for _, f := range transfer.Files {
		if f.FileIndex == fileIndex {
			return s.storage.CreateSignedURL(ctx, f.Bucket, f.StoragePath, 3600) // 1 hora
		}
	}
	return "", fmt.Errorf("archivo con índice %d no encontrado", fileIndex)
}

// sanitizeFilename retorna un nombre seguro, eliminando caracteres peligrosos.
func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	var safe strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			safe.WriteRune(r)
		} else {
			safe.WriteRune('_')
		}
	}
	if safe.Len() == 0 {
		return "archivo"
	}
	return safe.String()
}
