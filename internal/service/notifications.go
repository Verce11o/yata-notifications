package service

import (
	"github.com/Verce11o/yata-notifications/internal/repository"
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
