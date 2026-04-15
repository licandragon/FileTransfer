package models

import "time"

type Recipient struct {
	ID         string    `json:"id" db:"id"`
	TransferID string    `json:"transfer_id" db:"transfer_id"`
	Email      string    `json:"email" db:"email"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
