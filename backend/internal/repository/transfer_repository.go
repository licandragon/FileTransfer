package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/models"
)

type TransferRepository struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{db: db}
}

// Funcion para insertar un transfer en la base de datos
func (r *TransferRepository) Create(ctx context.Context, transfer *models.Transfer) error {

	query := `
	INSERT INTO transfers (id, download_token, user_id, expires_at)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		transfer.ID,
		transfer.DownloadToken,
		transfer.UserID,
		transfer.ExpiresAt,
	)

	return err
}

// Funcion que busca un transfer por su token
func (r *TransferRepository) GetByToken(ctx context.Context, token string) (*models.Transfer, error) {

	query := `
	SELECT id, download_token, user_id, expires_at, created_at
	FROM transfers
	WHERE download_token = $1
	`

	var transfer models.Transfer

	err := r.db.QueryRow(ctx, query, token).Scan(
		&transfer.ID,
		&transfer.DownloadToken,
		&transfer.UserID,
		&transfer.ExpiresAt,
		&transfer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transfer, nil
}
