package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/models"
)

type TransferRepository interface {
	Create(ctx context.Context, transfer *models.Transfer) error
	GetByToken(ctx context.Context, token string) (*models.Transfer, error)
	IncrementDownloadCount(ctx context.Context, token string) error
}
type transferRepo struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) TransferRepository {
	return &transferRepo{db: db}
}

// Create guarda el transfer y todos sus archivos en una sola transacción SQL.
func (r *transferRepo) Create(ctx context.Context, transfer *models.Transfer) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	// Si la función termina sin llegar al Commit, se hace Rollback automático.
	defer tx.Rollback(ctx)

	// 1. Insertar la cabecera de la transferencia
	queryTransfer := `
	INSERT INTO transfers (
        download_token, sender_email, subject_email, message_email,
        recipients, user_id, expires_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, download_count, created_at
	`

	err = tx.QueryRow(ctx, queryTransfer,
		transfer.DownloadToken,
		transfer.SenderEmail,
		transfer.SubjectEmail,
		transfer.MessageEmail,
		transfer.Recipients, // pgx enviará el []string como un array nativo de Postgres
		transfer.UserID,
		transfer.ExpiresAt,
	).Scan(&transfer.ID, &transfer.DownloadCount, &transfer.CreatedAt)

	if err != nil {
		return fmt.Errorf("error al insertar transfer: %w", err)
	}

	// 2. Insertar los archivos vinculados (Bulk Insert manual dentro de la TX)
	queryFile := `
	INSERT INTO files (transfer_id, filename, original_name, size, mime_type, storage_path, bucket)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at
	`

	for i := range transfer.Files {
		transfer.Files[i].TransferID = transfer.ID // Sincronizamos el ID recién generado

		err := tx.QueryRow(ctx, queryFile,
			transfer.Files[i].TransferID,
			transfer.Files[i].Filename,
			transfer.Files[i].OriginalName,
			transfer.Files[i].Size,
			transfer.Files[i].MimeType,
			transfer.Files[i].StoragePath,
			transfer.Files[i].Bucket,
		).Scan(&transfer.Files[i].ID, &transfer.Files[i].CreatedAt)

		if err != nil {
			return fmt.Errorf("error al insertar archivo %s: %w", transfer.Files[i].OriginalName, err)
		}
	}

	// 3. Confirmar todos los cambios
	return tx.Commit(ctx)
}

// GetByToken obtiene el transfer y carga automáticamente su lista de archivos.
func (r *transferRepo) GetByToken(ctx context.Context, token string) (*models.Transfer, error) {
	query := `
		SELECT id, download_token, download_count, sender_email, subject_email,
		message_email, recipients, user_id, expires_at, created_at
		FROM transfers
		WHERE download_token = $1
		`

	var t models.Transfer
	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID,
		&t.DownloadToken,
		&t.DownloadCount,
		&t.SenderEmail,
		&t.SubjectEmail,
		&t.MessageEmail,
		&t.Recipients,
		&t.UserID,
		&t.ExpiresAt,
		&t.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transfer no encontrado")
		}
		return nil, err
	}

	// Cargamos los archivos asociados usando la lógica que estaba en file_repository
	files, err := r.getFilesByTransferID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	t.Files = files

	return &t, nil
}

// getFilesByTransferID es un método auxiliar privado para mantener la limpieza.
func (r *transferRepo) getFilesByTransferID(ctx context.Context, transferID string) ([]models.File, error) {
	query := `
	SELECT id, transfer_id, filename, original_name, size, mime_type, storage_path, bucket, created_at
	FROM files
	WHERE transfer_id = $1
	`
	rows, err := r.db.Query(ctx, query, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		err := rows.Scan(
			&f.ID, &f.TransferID, &f.Filename, &f.OriginalName,
			&f.Size, &f.MimeType, &f.StoragePath, &f.Bucket, &f.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *transferRepo) IncrementDownloadCount(ctx context.Context, token string) error {
	query := `UPDATE transfers SET download_count = download_count + 1 WHERE download_token = $1`
	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no se encontró el token")
	}
	return nil
}
