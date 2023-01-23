package pg

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"github.com/SafeRE-IT/notifications-router-svc/internal/data"
)

const notificationsTableName = "notifications"

func NewNotificationsQ(db *pgdb.DB) data.NotificationsQ {
	return &notificationsQ{
		db:  db.Clone(),
		sql: sq.Select("n.*").From(fmt.Sprintf("%s as n", notificationsTableName)),
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

func (q *notificationsQ) Page(pageParams pgdb.OffsetPageParams) data.NotificationsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *notificationsQ) FilterByID(ids ...int64) data.NotificationsQ {
	q.sql = q.sql.Where(sq.Eq{"n.id": ids})
	return q
}

func (q *notificationsQ) FilterByToken(tokens ...string) data.NotificationsQ {
	q.sql = q.sql.Where(sq.Eq{"n.token": tokens})
	return q
}

func (q *notificationsQ) FilterByDestination(destination string, destinationType string) data.NotificationsQ {
	q.sql = q.sql.Join(fmt.Sprintf("%s as delivery on delivery.notification_id = n.id", deliveriesTableName)).
		Where(sq.Eq{"delivery.destination": destination}).
		Where(sq.Eq{"delivery.destination_type": destinationType})
	return q
}

func (q *notificationsQ) FilterByTopic(topics ...string) data.NotificationsQ {
	q.sql = q.sql.Where(sq.Eq{"n.topic": topics})
	return q
}

func (q *notificationsQ) FilterByScheduledAfter(time time.Time) data.NotificationsQ {
	q.sql = q.sql.Where(sq.GtOrEq{"scheduled_for": time})
	return q
}

func (q *notificationsQ) FilterByScheduledBefore(time time.Time) data.NotificationsQ {
	q.sql = q.sql.Where(sq.LtOrEq{"scheduled_for": time})
	return q
}
