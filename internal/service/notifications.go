package service

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NotificationsService struct {
	log    *zap.SugaredLogger
	tracer trace.Tracer
	repo   repository.Repository
}

func NewNotificationsService(log *zap.SugaredLogger, tracer trace.Tracer, repo repository.Repository) *NotificationsService {
	return &NotificationsService{log: log, tracer: tracer, repo: repo}
}

func (n NotificationsService) SubscribeToUser(ctx context.Context, request *pb.SubscribeToUserRequest) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.SubscribeToUser")
	defer span.End()

	return nil

}

func (n NotificationsService) UnSubscribeFromUser(ctx context.Context, request *pb.UnSubscribeFromUserRequest) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.UnSubscribeFromUser")
	defer span.End()

	return nil
}
