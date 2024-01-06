package repository

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/domain"
)

type RedisRepository interface {
	GetNotificationsByUserID(ctx context.Context, key string) ([]domain.Notification, error)
	SetNotificationsByUserID(ctx context.Context, key string, notifications []domain.Notification) error
	DeleteNotificationsByUserID(ctx context.Context, key string) error
}
