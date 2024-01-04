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

func (n *NotificationsPostgres) GetUserSubscription(ctx context.Context, userID string, toUserID string) (*domain.Subscriber, error) {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.SubscribeUser")
	defer span.End()

	q := "SELECT user_id, to_user_id, created_at, updated_at FROM subscribers WHERE user_id = $1 AND to_user_id = $2"

	var subscription domain.Subscriber

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

func (n *NotificationsPostgres) GetUserSubscribers(ctx context.Context, userID string) ([]domain.Subscriber, error) {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.GetUserSubscribers")
	defer span.End()

	q := "SELECT user_id, to_user_id, created_at, updated_at FROM subscribers WHERE to_user_id = $1"

	var result []domain.Subscriber

	err := sqlx.SelectContext(ctx, n.db, &result, q, userID)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (n *NotificationsPostgres) GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.GetNotifications")
	defer span.End()

	q := "SELECT * FROM notifications WHERE to_user_id = $1 ORDER BY read DESC, created_at DESC LIMIT 30"

	var result []domain.Notification

	err := sqlx.SelectContext(ctx, n.db, &result, q, userID)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func (n *NotificationsPostgres) MarkNotificationAsRead(ctx context.Context, userID string, notificationID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.GetNotifications")
	defer span.End()

	q := "UPDATE notifications SET read = TRUE WHERE to_user_id = $1 AND notification_id = $2"

	res, err := n.db.ExecContext(ctx, q, userID, notificationID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (n *NotificationsPostgres) ReadAllNotifications(ctx context.Context, userID string) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.ReadAllNotifications")
	defer span.End()

	q := "UPDATE notifications SET read = TRUE WHERE to_user_id = $1"

	res, err := n.db.ExecContext(ctx, q, userID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (n *NotificationsPostgres) BatchAddNotification(ctx context.Context, subscribers []domain.Subscriber, input domain.IncomingNewNotification) error {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.BatchAddNotification")
	defer span.End()

	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := "INSERT INTO notifications (to_user_id, from_user_id, type) VALUES ($1, $2, $3)"

	for _, sub := range subscribers {
		_, err = tx.ExecContext(ctx, q, sub.UserID, input.SenderID, input.Type)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (n *NotificationsPostgres) GetNotificationByID(ctx context.Context, userID string, notificationID string) (domain.Notification, error) {
	ctx, span := n.tracer.Start(ctx, "notificationsPostgres.GetNotificationByID")
	defer span.End()

	q := "SELECT * FROM notifications WHERE notification_id = $1 AND to_user_id = $2"

	var notification domain.Notification

	err := n.db.QueryRowxContext(ctx, q, notificationID, userID).StructScan(&notification)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Notification{}, sql.ErrNoRows
	}

	if err != nil {
		return domain.Notification{}, err
	}

	return notification, nil
}
