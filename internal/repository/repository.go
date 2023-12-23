package repository

import (
	"context"
)

type Repository interface {
	SubscribeToUser(ctx context.Context, userID, toUserID string) error
	UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error
}
