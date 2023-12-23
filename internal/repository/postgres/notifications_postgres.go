package postgres

import (
	"context"
	"database/sql"
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

func (n *NotificationsPostgres) SubscribeToUser(ctx context.Context, userId, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.SubscribeUser")
	defer span.End()

	q := "INSERT INTO subscribers(user_id, to_user_id) VALUES ($1, $2)"

	_, err := n.db.ExecContext(ctx, q, userId, toUserID)
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationsPostgres) UnSubscribeFromUser(ctx context.Context, userId, toUserID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.UnSubscribeFromUser")
	defer span.End()

	q := "DELETE FROM subscribers WHERE user_id = $1 AND to_user_id = $2"

	res, err := n.db.ExecContext(ctx, q, userId, toUserID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
