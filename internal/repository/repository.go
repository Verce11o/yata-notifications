package repository

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/domain"
)

type Repository interface {
	SubscribeToUser(ctx context.Context, userID, toUserID string) error
	GetUserSubscription(ctx context.Context, userID string, toUserID string) (*domain.UserSubscription, error)
	UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error
}
