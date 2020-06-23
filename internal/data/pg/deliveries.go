package pg

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
)

const deliveriesTableName = "deliveries"

func NewDeliveriesQ(db *pgdb.DB) data.DeliveriesQ {
	return &deliveriesQ{
		db:        db.Clone(),
		sql:       sq.Select("deliveries.*").From(deliveriesTableName),
		sqlUpdate: sq.Update(deliveriesTableName).Suffix("returning *"),
	}
}

type deliveriesQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *deliveriesQ) New() data.DeliveriesQ {
	return NewDeliveriesQ(q.db)
}

func (q *deliveriesQ) Get() (*data.Delivery, error) {
	var result data.Delivery
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *deliveriesQ) Select() ([]data.Delivery, error) {
	var result []data.Delivery
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *deliveriesQ) Update() ([]data.Delivery, error) {
	var result []data.Delivery
	err := q.db.Select(&result, q.sqlUpdate)

	return result, err
}

func (q *deliveriesQ) Transaction(fn func(q data.DeliveriesQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *deliveriesQ) FilterByNotificationID(ids ...int64) data.DeliveriesQ {
	stmt := sq.Eq{"deliveries.notification_id": ids}
	q.sql = q.sql.Where(stmt)
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) FilterByDestination(destinations ...string) data.DeliveriesQ {
	stmt := sq.Eq{"deliveries.destination": destinations}
	q.sql = q.sql.Where(stmt)
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) FilterByDestinationType(destinationTypes ...string) data.DeliveriesQ {
	stmt := sq.Eq{"deliveries.destination_type": destinationTypes}
	q.sql = q.sql.Where(stmt)
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) FilterByStatus(statuses ...data.DeliveryStatus) data.DeliveriesQ {
	stmt := sq.Eq{"deliveries.status": statuses}
	q.sql = q.sql.Where(stmt)
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) FilterByScheduledForAfter(time time.Time) data.DeliveriesQ {
	stmt := sq.GtOrEq{"notification.scheduled_for": time}
	q.sql = q.sql.Where(stmt)
	// TODO: Will not work for update
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) FilterById(ids ...int64) data.DeliveriesQ {
	stmt := sq.Eq{"deliveries.id": ids}
	q.sql = q.sql.Where(stmt)
	q.sqlUpdate = q.sqlUpdate.Where(stmt)
	return q
}

func (q *deliveriesQ) OrderByPriority(order string) data.DeliveriesQ {
	q.sql = q.sql.OrderBy(fmt.Sprintf("notification.priority %s", order))
	return q
}

// TODO: Find better way to join notification in different filters
func (q *deliveriesQ) JoinNotification() data.DeliveriesQ {
	stmt := fmt.Sprintf("%s as notification on notification.id = deliveries.id",
		notificationsTableName)
	q.sql = q.sql.Join(stmt)
	return q
}

func (q *deliveriesQ) SetStatus(status data.DeliveryStatus) data.DeliveriesQ {
	q.sqlUpdate = q.sqlUpdate.Set("status", status)
	return q
}
