package domain

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        uuid.UUID `json:"notification_id" db:"notification_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	SenderID  uuid.UUID `json:"sender_id,omitempty" db:"sender_id"`
	Read      bool      `json:"read" db:"read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type IncomingNewNotification struct {
	SenderID uuid.UUID `json:"sender_id"`
	Type     string    `json:"type"`
}
