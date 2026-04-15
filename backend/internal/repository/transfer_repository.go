package repository

import (
	"context"
	"encoding/json"
	"fmt"

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
	recipientsJSON, err := json.Marshal(transfer.Recipients)
	if err != nil {
		return fmt.Errorf("error al serializar destinatarios: %w", err)
	}

	query := `
	INSERT INTO transfers (
        download_token,
        sender_email,
        subject_email,
        message_email,
        recipients,
        user_id,
        expires_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, download_count, created_at
	`

	err = r.db.QueryRow(
		ctx,
		query,
		transfer.DownloadToken,
		transfer.SenderEmail,
		transfer.SubjectEmail,
		transfer.MessageEmail,
		recipientsJSON,
		transfer.UserID,
		transfer.ExpiresAt,
	).Scan(&transfer.ID, &transfer.DownloadCount, &transfer.CreatedAt)

	return err
}

// Funcion que busca un transfer por su token
func (r *TransferRepository) GetByToken(ctx context.Context, token string) (*models.Transfer, error) {

	query := `
	SELECT id, download_token, download_count, sender_email, subject_email,
	    message_email, recipients, user_id, expires_at, created_at
	FROM transfers
	WHERE download_token = $1
	`

	var transfer models.Transfer
	var recipientsData []byte

	err := r.db.QueryRow(ctx, query, token).Scan(
		&transfer.ID,
		&transfer.DownloadToken,
		&transfer.DownloadCount,
		&transfer.SenderEmail,
		&transfer.SubjectEmail,
		&transfer.MessageEmail,
		&recipientsData,
		&transfer.UserID,
		&transfer.ExpiresAt,
		&transfer.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(recipientsData, &transfer.Recipients); err != nil {
		return nil, fmt.Errorf("error al deserializar destinatarios: %w", err)
	}

	return &transfer, nil
}

func (r *TransferRepository) IncrementDownloadCount(ctx context.Context, token string) error {
	query := `
	UPDATE transfers
	SET download_count = download_count + 1
	WHERE download_token = $1
	`
	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("Error al incrementar contador de descargas: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No se encontro transfer con ese token: %s", token)
	}

	return nil
}
