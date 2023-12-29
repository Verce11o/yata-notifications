package grpc

import (
	"context"
	"github.com/Verce11o/yata-notifications/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-notifications/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type NotificationGRPC struct {
	log     *zap.SugaredLogger
	tracer  trace.Tracer
	service service.Notifications
	pb.UnimplementedNotificationsServer
}

func NewNotificationGRPC(log *zap.SugaredLogger, tracer trace.Tracer, service service.Notifications) *NotificationGRPC {
	return &NotificationGRPC{log: log, tracer: tracer, service: service}
}

func (n *NotificationGRPC) SubscribeToUser(ctx context.Context, input *pb.SubscribeToUserRequest) (*pb.SubscribeToUserResponse, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.SubscribeToUser")
	defer span.End()

	err := n.service.SubscribeToUser(ctx, input)

	if err != nil {
		n.log.Errorf("SubscribeToUser: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "SubscribeToUser: %v", err)
	}

	return &pb.SubscribeToUserResponse{}, nil
}

func (n *NotificationGRPC) UnSubscribeFromUser(ctx context.Context, input *pb.UnSubscribeFromUserRequest) (*pb.UnSubscribeFromUserResponse, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.UnSubscribeFromUser")
	defer span.End()

	err := n.service.UnSubscribeFromUser(ctx, input)

	if err != nil {
		n.log.Errorf("UnSubscribeFromUser: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "UnSubscribeFromUser: %v", err)
	}

	return &pb.UnSubscribeFromUserResponse{}, nil

}

func (n *NotificationGRPC) GetNotifications(ctx context.Context, input *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.GetNotifications")
	defer span.End()

	notifications, err := n.service.GetNotifications(ctx, input.GetUserId())

	if err != nil {
		n.log.Errorf("GetNotifications: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "UnSubscribeFromUser: %v", err)
	}

	return &pb.GetNotificationsResponse{Notifications: notifications}, nil

}
func (n *NotificationGRPC) MarkNotificationAsRead(ctx context.Context, input *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.MarkNotificationAsRead")
	defer span.End()

	err := n.service.MarkNotificationAsRead(ctx, input.GetUserId(), input.GetNotificationId())

	if err != nil {
		n.log.Errorf("MarkNotificationAsRead: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "MarkNotificationAsRead: %v", err)
	}

	return &pb.MarkNotificationAsReadResponse{}, nil
}
func (n *NotificationGRPC) ReadAllNotifications(ctx context.Context, input *pb.ReadAllNotificationsRequest) (*pb.ReadAllNotificationsResponse, error) {
	ctx, span := n.tracer.Start(ctx, "notificationService.ReadAllNotifications")
	defer span.End()

	err := n.service.ReadAllNotifications(ctx, input.GetUserId())

	if err != nil {
		n.log.Errorf("ReadAllNotifications: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "ReadAllNotifications: %v", err)
	}

	return &pb.ReadAllNotificationsResponse{}, nil

}
