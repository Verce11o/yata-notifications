package service

import (
	"context"
	pb "github.com/Verce11o/yata-protos/gen/go/notifications"
)

type Notifications interface {
	SubscribeToUser(ctx context.Context, request *pb.SubscribeToUserRequest) error
	UnSubscribeFromUser(ctx context.Context, request *pb.UnSubscribeFromUserRequest) error
}
