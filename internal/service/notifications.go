package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Verce11o/yata-notifications/internal/domain"
	"github.com/Verce11o/yata-notifications/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-notifications/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NotificationsService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	repo   repository.Repository
	redis  repository.RedisRepository
}

func NewNotificationsService(log *zap.SugaredLogger, tracer trace.Tracer, repo repository.Repository, redis repository.RedisRepository) *NotificationsService {
	return &NotificationsService{log: log, tracer: tracer, repo: repo, redis: redis}
}

func (n *NotificationsService) SubscribeToUser(ctx context.Context, request *pb.SubscribeToUserRequest) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.SubscribeToUser")
	defer span.End()

	subscription, err := n.repo.GetUserSubscription(ctx, request.GetUserId(), request.GetToUserId())
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		n.log.Errorf("cannot get user subscription by id %v", err.Error())
		return err
	}

	if subscription != nil {
		n.log.Infof("user already subscribed")
		return grpc_errors.ErrSubAlreadyExists
	}

	if request.GetToUserId() == request.GetUserId() {
		n.log.Errorf("user cannot subscribe to himself")
		return grpc_errors.ErrInvalidUser
	}

	err = n.repo.SubscribeToUser(ctx, request.GetUserId(), request.GetToUserId())

	if err != nil {
		n.log.Errorf("cannot subscribe user: %v", err.Error())
		return err
	}

	return nil

}

func (n *NotificationsService) UnSubscribeFromUser(ctx context.Context, request *pb.UnSubscribeFromUserRequest) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.UnSubscribeFromUser")
	defer span.End()

	err := n.repo.UnSubscribeFromUser(ctx, request.GetUserId(), request.GetToUserId())

	if err != nil {
		n.log.Errorf("cannot unsubscribe user: %v", err.Error())
		return err
	}

	return nil
}

func (n *NotificationsService) GetNotifications(ctx context.Context, userID string) ([]*pb.Notification, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.GetNotifications")
	defer span.End()

	cachedNotifications, err := n.redis.GetNotificationsByUserID(ctx, userID)
	if err != nil {
		n.log.Errorf("cannot get cached notifications: %v", err.Error())
	}

	if cachedNotifications != nil {
		return domainToNotificationPb(cachedNotifications), nil
	}

	notifications, err := n.repo.GetNotifications(ctx, userID)

	if err != nil {
		n.log.Errorf("cannot get notifications: %v", err.Error())
		return nil, err
	}

	if err := n.redis.SetNotificationsByUserID(ctx, userID, notifications); err != nil {
		n.log.Errorf("cannot set notifications in redis: %v", err.Error())
	}

	return domainToNotificationPb(notifications), nil
}

func (n *NotificationsService) MarkNotificationAsRead(ctx context.Context, userID string, notificationID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.MarkNotificationAsRead")
	defer span.End()

	notification, err := n.repo.GetNotificationByID(ctx, userID, notificationID)
	if err != nil {
		n.log.Errorf("cannot get notification by id: %v", err)
		return err
	}

	if notification.ToUserID.String() != userID {
		return grpc_errors.ErrPermissionDenied
	}

	err = n.repo.MarkNotificationAsRead(ctx, userID, notificationID)

	if err != nil {
		n.log.Errorf("cannot mark notification as read: %v", err)
		return err
	}

	err = n.redis.DeleteNotificationsByUserID(ctx, userID)

	if err != nil {
		n.log.Errorf("cannot delete notificaion in redis")
		return err
	}

	return err
}

func (n *NotificationsService) ReadAllNotifications(ctx context.Context, userID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.ReadAllNotifications")
	defer span.End()

	err := n.repo.ReadAllNotifications(ctx, userID)

	if err != nil {
		n.log.Errorf("cannot read all notifications: %v", err.Error())
		return err
	}

	err = n.redis.DeleteNotificationsByUserID(ctx, userID)
	if err != nil {
		n.log.Errorf("cannot delete notificaion in redis")
		return err
	}
	return nil
}

func (n *NotificationsService) GetUserSubscribers(ctx context.Context, userID string) ([]domain.Subscriber, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.GetUserSubscribers")
	defer span.End()

	subs, err := n.repo.GetUserSubscribers(ctx, userID)
	if err != nil {
		n.log.Errorf("cannot get user subscribers: %v", err.Error())
		return nil, err
	}

	return subs, nil
}

func (n *NotificationsService) GetUserSubscriptions(ctx context.Context, userID string, cursor string) ([]*pb.Subscriber, string, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.GetUserSubscriptions")
	defer span.End()

	subscriptions, cursor, err := n.repo.GetUserSubscriptions(ctx, userID, cursor)
	if err != nil {
		n.log.Errorf("cannot get user subscriptions: %v", err.Error())
		return nil, "", err
	}

	return domainToSubscriberPb(subscriptions), cursor, nil
}

func (n *NotificationsService) BatchAddNotification(ctx context.Context, subscribers []domain.Subscriber, notification domain.IncomingNewNotification) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.BatchAddNotification")
	defer span.End()

	err := n.repo.BatchAddNotification(ctx, subscribers, notification)

	if err != nil {
		n.log.Errorf("cannot add new notification: %v", err.Error())
		return err
	}

	for _, sub := range subscribers {
		if err := n.redis.DeleteNotificationsByUserID(ctx, sub.UserID); err != nil {
			n.log.Errorf("cannot delete user notification cache: %v", err.Error())
			return err
		}
	}

	return nil

}

func domainToNotificationPb(notifications []domain.Notification) []*pb.Notification {
	result := make([]*pb.Notification, 0, len(notifications))

	for _, notification := range notifications {
		result = append(result, &pb.Notification{
			NotificationId: notification.NotificationID.String(),
			UserId:         notification.ToUserID.String(),
			SenderId:       notification.FromUserID.String(),
			Read:           notification.Read,
			CreatedAt:      timestamppb.New(notification.CreatedAt),
			Type:           notification.Type,
		})
	}
	return result
}

func domainToSubscriberPb(subscribers []domain.Subscriber) []*pb.Subscriber {
	result := make([]*pb.Subscriber, 0, len(subscribers))

	for _, subscriber := range subscribers {
		result = append(result, &pb.Subscriber{
			UserId:    subscriber.UserID,
			ToUserId:  subscriber.ToUserID,
			CreatedAt: timestamppb.New(subscriber.CreatedAt),
		})
	}
	return result
}
