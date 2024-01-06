package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Verce11o/yata-notifications/internal/domain"
	"github.com/Verce11o/yata-notifications/internal/lib/grpc_errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	notificationTTL = 3600
)

type NotificationRedis struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewNotificationRedis(client *redis.Client, tracer trace.Tracer) *NotificationRedis {
	return &NotificationRedis{client: client, tracer: tracer}
}

func (n *NotificationRedis) GetNotificationsByUserID(ctx context.Context, key string) ([]domain.Notification, error) {
	ctx, span := n.tracer.Start(ctx, "notificationRedis.GetNotificationsByID")
	defer span.End()

	notificationBytes, err := n.client.Get(ctx, n.createKey(key)).Bytes()

	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, grpc_errors.ErrNotFound
		}
		return nil, err
	}

	var notifications []domain.Notification
	if err := json.Unmarshal(notificationBytes, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil

}

func (n *NotificationRedis) SetNotificationsByUserID(ctx context.Context, key string, notifications []domain.Notification) error {
	ctx, span := n.tracer.Start(ctx, "notificationRedis.SetNotificationsByUserID")
	defer span.End()

	notificationBytes, err := json.Marshal(notifications)

	if err != nil {
		return err
	}

	return n.client.Set(ctx, n.createKey(key), notificationBytes, time.Second*time.Duration(notificationTTL)).Err()

}

func (n *NotificationRedis) DeleteNotificationsByUserID(ctx context.Context, key string) error {
	ctx, span := n.tracer.Start(ctx, "notificationRedis.DeleteNotificationsByUserID")
	defer span.End()

	return n.client.Del(ctx, n.createKey(key)).Err()

}

func (n *NotificationRedis) createKey(key string) string {
	return fmt.Sprintf("notification:%s", key)
}
