package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/models"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, transfer *models.Transfer) error

	AddOrRetryFile(ctx context.Context, file *models.File) (bool, error)
	CompleteTransfer(ctx context.Context, transferID, uploadToken uuid.UUID) (uuid.UUID, error)
	GetByUploadToken(ctx context.Context, uploadToken uuid.UUID) (*models.Transfer, error)
	GetByDownloadToken(ctx context.Context, downloadToken uuid.UUID) (*models.Transfer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transfer, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Transfer, error)
	IncrementDownloadCount(ctx context.Context, downloadToken uuid.UUID) error
	DeletePendingTransfer(ctx context.Context, transferID uuid.UUID) error
	CountUploadedFiles(ctx context.Context, transferID uuid.UUID) (int, error)
}
type transferRepo struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) TransferRepository {
	return &transferRepo{db: db}
}

// Create guarda el transfer y todos sus archivos en una sola transacción SQL.
func (r *transferRepo) CreateTransfer(ctx context.Context, transfer *models.Transfer) error {
	transfer.UploadToken = uuid.New()
	transfer.DownloadToken = uuid.New()
	transfer.StatusTransfer = "pending"
	fmt.Println(transfer)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %w", err)
	}
	// Si la función termina sin llegar al Commit, se hace Rollback automático.
	defer tx.Rollback(ctx)

	query := `
    INSERT INTO transfers (
        download_token, upload_token, sender_email, subject_email,
        message_email, recipients, user_id, expires_at, status_transfer, total_files
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id, download_count, created_at`

	err = tx.QueryRow(ctx, query,
		transfer.DownloadToken,
		transfer.UploadToken,
		transfer.SenderEmail,
		transfer.SubjectEmail,
		transfer.MessageEmail,
		transfer.Recipients,
		transfer.UserID,
		transfer.ExpiresAt,
		transfer.StatusTransfer,
		transfer.TotalFiles,
	).Scan(&transfer.ID, &transfer.DownloadCount, &transfer.CreatedAt)

	if err != nil {
		return fmt.Errorf("error al crear transferencia: %w", err)
		// ⚠️ Aquí se retorna, y el defer ejecuta Rollback. Correcto.
	}

	// Confirmar todos los cambios
	return tx.Commit(ctx)
}

func (r *transferRepo) IncrementDownloadCount(ctx context.Context, token uuid.UUID) error {

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

// AddOrRetryFile inserta un archivo en un slot o lo reintenta si había fallado.
// Retorna true si el archivo era nuevo, false si fue un reintento.
// Devuelve error si el slot ya está en estado 'uploaded'.
func (r *transferRepo) AddOrRetryFile(ctx context.Context, file *models.File) (bool, error) {
	// Verificamos estado actual del slot
	var currentStatus string
	err := r.db.QueryRow(ctx,
		`SELECT status_file FROM files WHERE transfer_id = $1 AND file_index = $2`,
		file.TransferID, file.FileIndex,
	).Scan(&currentStatus)

	if err != nil && err != pgx.ErrNoRows {
		return false, fmt.Errorf("error consultando archivo: %w", err)
	}

	if err == nil { // ya existe
		if currentStatus == "uploaded" {
			return false, fmt.Errorf("el archivo %d ya fue subido exitosamente", file.FileIndex)
		}
		// Reintento de uno fallido: actualizamos metadatos
		queryUpdate := `
        UPDATE files SET
            file_name = $1, original_name = $2, size = $3,
            mime_type = $4, storage_path = $5, bucket = $6,
            status_file = 'uploaded', created_at = NOW()
        WHERE transfer_id = $7 AND file_index = $8
        RETURNING id, created_at`
		err = r.db.QueryRow(ctx, queryUpdate,
			file.Filename, file.OriginalName, file.SizeFile,
			file.MimeType, file.StoragePath, file.Bucket,
			file.TransferID, file.FileIndex,
		).Scan(&file.ID, &file.CreatedAt)
		if err != nil {
			return false, fmt.Errorf("error al reintentar archivo: %w", err)
		}
		return false, nil
	}

	// No existe, insertamos nuevo
	file.StatusFile = "uploaded"
	queryInsert := `
    INSERT INTO files (
        transfer_id, file_index, file_name, original_name,
        size_file, mime_type, storage_path, bucket, status_file
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id, created_at`
	err = r.db.QueryRow(ctx, queryInsert,
		file.TransferID, file.FileIndex, file.Filename,
		file.OriginalName, file.SizeFile, file.MimeType,
		file.StoragePath, file.Bucket, file.StatusFile,
	).Scan(&file.ID, &file.CreatedAt)
	if err != nil {
		return false, fmt.Errorf("error al insertar archivo: %w", err)
	}
	return true, nil
}

// CompleteTransfer verifica que el upload_token sea correcto, que todos los archivos
// estén subidos y actualiza el estado a 'complete'. Retorna el download_token.
func (r *transferRepo) CompleteTransfer(ctx context.Context, transferID, uploadToken uuid.UUID) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// Obtener la transferencia con bloqueo para evitar condiciones de carrera
	var transfer models.Transfer
	err = tx.QueryRow(ctx,
		`SELECT id, upload_token, download_token, total_files, status_transfer
        	FROM transfers WHERE id = $1 FOR UPDATE`, transferID,
	).Scan(&transfer.ID, &transfer.UploadToken, &transfer.DownloadToken,
		&transfer.TotalFiles, &transfer.StatusTransfer)
	if err != nil {
		return uuid.Nil, fmt.Errorf("transferencia no encontrada: %w", err)
	}

	if transfer.UploadToken != uploadToken {
		return uuid.Nil, fmt.Errorf("upload_token inválido")
	}
	if transfer.StatusTransfer != "pending" {
		return uuid.Nil, fmt.Errorf("la transferencia no está en estado pendiente")
	}

	// Contar archivos subidos correctamente
	var uploaded int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM files WHERE transfer_id = $1 AND status_file = 'uploaded'`,
		transferID,
	).Scan(&uploaded)
	if err != nil {
		return uuid.Nil, err
	}
	if uploaded != transfer.TotalFiles {
		return uuid.Nil, fmt.Errorf("faltan archivos: %d de %d subidos", uploaded, transfer.TotalFiles)
	}

	// Actualizar estado
	_, err = tx.Exec(ctx,
		`UPDATE transfers SET status_transfer = 'complete' WHERE id = $1`, transferID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error al completar transferencia: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return transfer.DownloadToken, nil
}

// CountUploadedFiles cuenta los archivos con status_file='uploaded' para un transfer.
func (r *transferRepo) CountUploadedFiles(ctx context.Context, transferID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM files WHERE transfer_id = $1 AND status_file = 'uploaded'`,
		transferID,
	).Scan(&count)
	return count, err
}

// DeletePendingTransfer elimina una transferencia pendiente y sus archivos (útil para limpieza).
func (r *transferRepo) DeletePendingTransfer(ctx context.Context, transferID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM transfers WHERE id = $1 AND status_transfer = 'pending'`,
		transferID,
	)
	return err
}

// GetByDownloadToken obtiene una transferencia pública para descarga.
func (r *transferRepo) GetByDownloadToken(ctx context.Context, downloadToken uuid.UUID) (*models.Transfer, error) {
	query := `
    SELECT id, download_token, upload_token, download_count, sender_email,
           subject_email, message_email, recipients, user_id, expires_at,
           status_transfer, total_files, created_at
    FROM transfers WHERE download_token = $1 AND status_transfer = 'complete'`
	var t models.Transfer
	err := r.db.QueryRow(ctx, query, downloadToken).Scan(
		&t.ID, &t.DownloadToken, &t.UploadToken, &t.DownloadCount,
		&t.SenderEmail, &t.SubjectEmail, &t.MessageEmail,
		&t.Recipients, &t.UserID, &t.ExpiresAt, &t.StatusTransfer, &t.TotalFiles, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transferencia no encontrada o no disponible")
		}
		return nil, err
	}
	files, err := r.getFilesByTransferID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	t.Files = files
	return &t, nil
}

// GetByID obtiene una transferencia por su UUID interno.
func (r *transferRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Transfer, error) {
	query := `
    SELECT id, download_token, upload_token, download_count, sender_email,
		subject_email, message_email, recipients, user_id, expires_at,
		status_transfer, total_files, created_at
    FROM transfers WHERE id = $1`
	var t models.Transfer
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.DownloadToken, &t.UploadToken, &t.DownloadCount,
		&t.SenderEmail, &t.SubjectEmail, &t.MessageEmail,
		&t.Recipients, &t.UserID, &t.ExpiresAt, &t.StatusTransfer, &t.TotalFiles, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transferencia no encontrada")
		}
		return nil, err
	}
	files, err := r.getFilesByTransferID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	t.Files = files
	return &t, nil
}

// GetByUploadToken obtiene una transferencia por su UUID interno.
func (r *transferRepo) GetByUploadToken(ctx context.Context, uploadToken uuid.UUID) (*models.Transfer, error) {
	query := `
    SELECT id, download_token, upload_token, download_count, sender_email,
		subject_email, message_email, recipients, user_id, expires_at,
		status_transfer, total_files, created_at
    FROM transfers WHERE upload_token = $1`
	var t models.Transfer
	err := r.db.QueryRow(ctx, query, uploadToken).Scan(
		&t.ID, &t.DownloadToken, &t.UploadToken, &t.DownloadCount,
		&t.SenderEmail, &t.SubjectEmail, &t.MessageEmail,
		&t.Recipients, &t.UserID, &t.ExpiresAt, &t.StatusTransfer, &t.TotalFiles, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transferencia no encontrada")
		}
		return nil, err
	}
	// Cargar archivos asociados
	files, err := r.getFilesByTransferID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	t.Files = files
	return &t, nil
}

// getFilesByTransferID es un helper privado que carga los archivos de una transferencia.
func (r *transferRepo) getFilesByTransferID(ctx context.Context, transferID uuid.UUID) ([]models.File, error) {
	query := `
    SELECT id, transfer_id, file_index, file_name, original_name, size_file,
           mime_type, storage_path, bucket, status_file, created_at
    FROM files
    WHERE transfer_id = $1
    ORDER BY file_index ASC`
	rows, err := r.db.Query(ctx, query, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var f models.File
		err := rows.Scan(
			&f.ID, &f.TransferID, &f.FileIndex, &f.Filename, &f.OriginalName,
			&f.SizeFile, &f.MimeType, &f.StoragePath, &f.Bucket, &f.StatusFile, &f.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

// ListByUser implements [TransferRepository]. Pendiente para implementar
func (r *transferRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Transfer, error) {
	panic("unimplemented")
}
