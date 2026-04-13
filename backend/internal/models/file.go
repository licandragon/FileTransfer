package models

import "time"

type File struct {
	ID           string    `json:"id" db:"id"`
	TransferID   string    `json:"transfer_id" db:"transfer_id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_name"`
	Size         int64     `json:"size" db:"size"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	StoragePath  string    `json:"storage_path" db:"storage_path"`
	Bucket       string    `json:"bucket" db:"bucket"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
