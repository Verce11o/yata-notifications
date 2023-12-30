package domain

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	NotificationID uuid.UUID `json:"notification_id" db:"notification_id"`
	ToUserID       uuid.UUID `json:"to_user_id" db:"to_user_id"`
	FromUserID     uuid.UUID `json:"from_user_id,omitempty" db:"from_user_id"`
	Type           string    `json:"type" db:"type"`
	Read           bool      `json:"read" db:"read"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type IncomingNewNotification struct {
	SenderID uuid.UUID `json:"sender_id"`
	Type     string    `json:"type"`
}
