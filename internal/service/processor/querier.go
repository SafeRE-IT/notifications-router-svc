package processor

import (
	"errors"
	"time"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/kit/pgdb"
	"github.com/SafeRE-IT/notifications-router-svc/internal/data"
	"github.com/SafeRE-IT/notifications-router-svc/internal/data/pg"
)

func newQuerier(db *pgdb.DB) *querier {
	return &querier{
		deliveriesQ:    pg.NewDeliveriesQ(db),
		notificationsQ: pg.NewNotificationsQ(db),
	}
}

type querier struct {
	deliveriesQ    data.DeliveriesQ
	notificationsQ data.NotificationsQ
}

func (q *querier) getPendingDeliveries() ([]data.Delivery, error) {
	return q.deliveriesQ.New().
		JoinNotification().
		FilterByStatus(resources.DeliveryStatusNotSent).
		FilterByScheduledBefore(time.Now().UTC()).
		OrderByPriority(pgdb.OrderTypeDesc).
		Select()
}

// TODO: Join into deliveries
func (q *querier) getNotification(id int64) (data.Notification, error) {
	result, err := q.notificationsQ.New().
		FilterByID(id).
		Get()
	if result == nil {
		return data.Notification{}, errors.New("failed to find notification")
	}
	return *result, err
}

func (q *querier) setDeliveryStatus(id int64, status resources.DeliveryStatus) error {
	_, err := q.deliveriesQ.New().
		FilterById(id).
		SetStatus(status).
		Update()
	return err
}
