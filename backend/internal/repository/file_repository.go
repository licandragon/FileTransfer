package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/licandragon/FileTransfer/backend/internal/models"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

// Funcion para insertar la meta data del archivo
func (r *FileRepository) Create(ctx context.Context, file *models.File) error {

	query := `
	INSERT INTO files
	(id, transfer_id, filename, original_name, size, mime_type, storage_path, bucket)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		file.ID,
		file.TransferID,
		file.Filename,
		file.OriginalName,
		file.Size,
		file.MimeType,
		file.StoragePath,
		file.Bucket,
	)

	return err
}

// Funcion que obtiene los archivos por el transferID
func (r *FileRepository) GetByTransferID(ctx context.Context, transferID string) ([]models.File, error) {

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

		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.TransferID,
			&file.Filename,
			&file.OriginalName,
			&file.Size,
			&file.MimeType,
			&file.StoragePath,
			&file.Bucket,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
