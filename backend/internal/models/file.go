package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TransferID   uuid.UUID `json:"transfer_id" db:"transfer_id"`
	Filename     string    `json:"file_name" db:"file_name"`
	OriginalName string    `json:"original_name" db:"original_name"`
	SizeFile     int64     `json:"siz_filee" db:"size_file"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	StoragePath  string    `json:"storage_path" db:"storage_path"`
	FileIndex    int       `json:"file_index" db:"file_index"`
	Bucket       string    `json:"bucket" db:"bucket"`
	StatusFile   string    `json:"status_file" db:"status_file"` // "uploaded" o "failed"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
