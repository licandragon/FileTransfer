package models

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	DownloadToken  uuid.UUID  `json:"download_token" db:"download_token"`       //token para los enlaces de descarga
	UploadToken    uuid.UUID  `json:"upload_token,omitempty" db:"upload_token"` //Token para la subida de los archivos
	SenderEmail    string     `json:"sender_email" db:"sender_email"`
	SubjectEmail   string     `json:"subject_email" db:"subject_email"`
	MessageEmail   string     `json:"message_email" db:"message_email"`
	Recipients     []string   `json:"recipients" db:"recipients"`
	UserID         *uuid.UUID `json:"user_id" db:"user_id"`
	StatusTransfer string     `json:"status_transfer" db:"status_transfer"`
	DownloadCount  int        `json:"download_count" db:"download_count"`
	TotalFiles     int        `json:"total_files" db:"total_files"`
	ExpiresAt      *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`

	Files []File `json:"files,omitempty"`
}
