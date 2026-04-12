package models

import "time"

type File struct {
	ID        int       `json:"id" db:"id"`
	Filename  string    `json:"filename" db:"filename"`
	Token     string    `json:"token" db:"token"`
	Size      int64     `json:"size" db:"size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
