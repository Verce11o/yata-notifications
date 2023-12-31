package grpc_errors

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
)

var (
	ErrAddMinio         = errors.New("add file error")
	ErrNotFound         = errors.New("not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidCursor    = errors.New("invalid pagination cursor")
	ErrSubAlreadyExists = errors.New("already subscribed")
	ErrInvalidUser      = errors.New("invalid user")
)

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrAddMinio):
		return codes.Internal
	case errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, ErrPermissionDenied):
		return codes.PermissionDenied
	case errors.Is(err, redis.Nil):
		return codes.NotFound
	case errors.Is(err, ErrSubAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrInvalidUser):
		return codes.InvalidArgument
	case errors.Is(err, ErrInvalidCursor):
		return codes.InvalidArgument
	}
	return codes.Internal
}
