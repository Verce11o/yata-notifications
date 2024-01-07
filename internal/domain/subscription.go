package domain

import "time"

type Subscriber struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userID" db:"user_id"`
	ToUserID  string    `json:"toUserID" db:"to_user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
