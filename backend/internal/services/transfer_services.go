package services

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/licandragon/FileTransfer/backend/internal/models"
	"github.com/licandragon/FileTransfer/backend/internal/repository"
	"github.com/licandragon/FileTransfer/backend/internal/storage"
)

// Se define interfaz del servicio
type TransferService interface {
	CreateTransfer(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error)
	AddFileToTransfer(ctx context.Context, fileHeader *multipart.FileHeader, uploadToken uuid.UUID, fileIndex int) (*models.File, error)
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

// AddFileToTransfer implements [TransferService].
func (t *transferService) AddFileToTransfer(ctx context.Context, fileHeader *multipart.FileHeader, uploadToken uuid.UUID, fileIndex int) (*models.File, error) {
	panic("unimplemented")
}

// CompleteTransfer implements [TransferService].
func (t *transferService) CompleteTransfer(ctx context.Context, uploadToken uuid.UUID) (uuid.UUID, error) {
	panic("unimplemented")
}

// CreateTransfer implements [TransferService].
func (t *transferService) CreateTransfer(ctx context.Context, transfer *models.Transfer) (*models.Transfer, error) {
	panic("unimplemented")
}

// DeleteTransfer implements [TransferService].
func (t *transferService) DeleteTransfer(ctx context.Context, transferID uuid.UUID, userID uuid.UUID) error {
	panic("unimplemented")
}

// GetFileSignedURL implements [TransferService].
func (t *transferService) GetFileSignedURL(ctx context.Context, downloadToken uuid.UUID, fileIndex int) (string, error) {
	panic("unimplemented")
}

// GetTransferByDownloadToken implements [TransferService].
func (t *transferService) GetTransferByDownloadToken(ctx context.Context, downloadToken uuid.UUID) (*models.Transfer, error) {
	panic("unimplemented")
}

// GetTransferByID implements [TransferService].
func (t *transferService) GetTransferByID(ctx context.Context, id uuid.UUID) (*models.Transfer, error) {
	panic("unimplemented")
}

// GetTransferByUploadToken implements [TransferService].
func (t *transferService) GetTransferByUploadToken(ctx context.Context, uploadToken uuid.UUID) (*models.Transfer, error) {
	panic("unimplemented")
}

// ListUserTransfers implements [TransferService].
func (t *transferService) ListUserTransfers(ctx context.Context, userID uuid.UUID) ([]models.Transfer, error) {
	panic("unimplemented")
}
