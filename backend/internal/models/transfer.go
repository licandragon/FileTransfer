package models

import "time"

type Transfer struct {
	ID            string     `json:"id" db:"id"`
	DownloadToken string     `json:"download_token" db:"download_token"`
	UserID        *string    `json:"user_id" db:"user_id"`
	ExpiresAt     *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`

	Files []File `json:"files,omitempty"`
}
