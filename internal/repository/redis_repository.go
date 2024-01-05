package repository

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/domain"
)

type RedisRepository interface {
	GetNotificationsByID(ctx context.Context, key string) ([]domain.Notification, error)
	SetNotificationsByID(ctx context.Context, key string, notification domain.Notification) error
	DeleteNotificationsByID(ctx context.Context, key string) error
}
