package repository

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/domain"
)

type Repository interface {
	SubscribeToUser(ctx context.Context, userID, toUserID string) error
	GetUserSubscription(ctx context.Context, userID string, toUserID string) (*domain.Subscriber, error)
	UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error
	GetUserSubscribers(ctx context.Context, userID string) ([]domain.Subscriber, error)

	BatchAddNotification(ctx context.Context, subscribers []domain.Subscriber, input domain.IncomingNewNotification) error
	GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID string, notificationID string) error
	ReadAllNotifications(ctx context.Context, userID string) error
}
