package pg

import (
	"database/sql"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
)

const notificationsTableName = "notifications"

func NewNotificationsQ(db *pgdb.DB) data.NotificationsQ {
	return &notificationsQ{
		db:  db.Clone(),
		sql: sq.Select("*").From(notificationsTableName),
	}
}

type notificationsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *notificationsQ) New() data.NotificationsQ {
	return NewNotificationsQ(q.db)
}

func (q *notificationsQ) Get() (*data.Notification, error) {
	var result data.Notification
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *notificationsQ) Select() ([]data.Notification, error) {
	var result []data.Notification
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *notificationsQ) Transaction(fn func(q data.NotificationsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *notificationsQ) Insert(value data.Notification) (data.Notification, error) {
	value.CreatedAt = time.Now().UTC()
	clauses := structs.Map(value)
	clauses["message"] = value.Message

	var result data.Notification
	stmt := sq.Insert(notificationsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *notificationsQ) InsertDeliveries(deliveries []data.Delivery) ([]data.Delivery, error) {
	if len(deliveries) == 0 {
		return nil, errors.New("empty array is not allowed")
	}

	names := []string{
		"notification_id",
		"destination",
		"destination_type",
		"status",
		"sent_at",
	}
	stmt := sq.Insert(deliveriesTableName).Columns(names...)
	for _, item := range deliveries {
		stmt = stmt.Values([]interface{}{
			item.NotificationID,
			item.Destination,
			item.DestinationType,
			item.Status,
			item.SentAt,
		}...)
	}

	stmt = stmt.Suffix("returning *")
	var result []data.Delivery
	err := q.db.Select(&result, stmt)

	return result, err
}
