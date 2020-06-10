package pg

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
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
