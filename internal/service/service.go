package service

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/domain"
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
)

type Notifications interface {
	SubscribeToUser(ctx context.Context, request *pb.SubscribeToUserRequest) error
	UnSubscribeFromUser(ctx context.Context, request *pb.UnSubscribeFromUserRequest) error
	GetUserSubscribers(ctx context.Context, userID string) ([]domain.Subscriber, error)
	BatchAddNotification(ctx context.Context, subscribers []domain.Subscriber, notification domain.IncomingNewNotification) error
	GetNotifications(ctx context.Context, userID string) ([]*pb.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID string, notificationID string) error
	ReadAllNotifications(ctx context.Context, userID string) error
}
