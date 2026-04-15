package models

import "time"

type Transfer struct {
	ID            string     `json:"id" db:"id"`
	DownloadToken string     `json:"download_token" db:"download_token"`
	DownloadCount int        `json:"download_count" db:"download_count"`
	SenderEmail   string     `json:"sender_email" db:"sender_email"`
	SubjectEmail  string     `json:"subject_email" db:"subject_email"`
	Message       *string    `json:"message" db:"message"`
	Recipients    any        `json:"recipients" db:"recipients"`
	UserID        *string    `json:"user_id" db:"user_id"`
	ExpiresAt     *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`

	Files []File `json:"files,omitempty"`
}
