package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Verce11o/yata-notifications/internal/domain"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type NotificationsPostgres struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewNotificationsPostgres(db *sqlx.DB, tracer trace.Tracer) *NotificationsPostgres {
	return &NotificationsPostgres{db: db, tracer: tracer}
}

func (n *NotificationsPostgres) SubscribeToUser(ctx context.Context, userID, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.SubscribeUser")
	defer span.End()

	q := "INSERT INTO subscribers(user_id, to_user_id) VALUES ($1, $2)"

	_, err := n.db.ExecContext(ctx, q, userID, toUserID)
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationsPostgres) GetUserSubscription(ctx context.Context, userID string, toUserID string) (*domain.UserSubscription, error) {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.SubscribeUser")
	defer span.End()

	q := "SELECT user_id, to_user_id, created_at, updated_at FROM subscribers WHERE user_id = $1 AND to_user_id = $2"

	var subscription domain.UserSubscription

	err := n.db.QueryRowxContext(ctx, q, userID, toUserID).StructScan(&subscription)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	return &subscription, err

}

func (n *NotificationsPostgres) UnSubscribeFromUser(ctx context.Context, userID, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.UnSubscribeFromUser")
	defer span.End()

	q := "DELETE FROM subscribers WHERE user_id = $1 AND to_user_id = $2"

	res, err := n.db.ExecContext(ctx, q, userID, toUserID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
