package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Verce11o/yata-notifications/internal/lib/grpc_errors"
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

func (n *NotificationsService) SubscribeToUser(ctx context.Context, request *pb.SubscribeToUserRequest) error {
	ctx, span := n.tracer.Start(ctx, "notificationService.SubscribeToUser")
	defer span.End()

	subscription, err := n.repo.GetUserSubscription(ctx, request.GetUserId(), request.GetToUserId())
	// todo check to_user_id for existence
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
