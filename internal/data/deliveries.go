package data

import (
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"
)

type DeliveriesQ interface {
	New() DeliveriesQ

	Get() (*Delivery, error)
	Select() ([]Delivery, error)
	Update() ([]Delivery, error)

	Transaction(fn func(q DeliveriesQ) error) error

	FilterByNotificationID(ids ...int64) DeliveriesQ
	FilterByDestination(destinations ...string) DeliveriesQ
	FilterByDestinationType(destinationTypes ...string) DeliveriesQ
	FilterByStatus(statuses ...resources.DeliveryStatus) DeliveriesQ
	FilterByScheduledBefore(time time.Time) DeliveriesQ
	FilterById(ids ...int64) DeliveriesQ

	OrderByPriority(order string) DeliveriesQ

	JoinNotification() DeliveriesQ

	SetStatus(status resources.DeliveryStatus) DeliveriesQ
}

type Delivery struct {
	ID              int64                    `db:"id" structs:"-"`
	NotificationID  int64                    `db:"notification_id" structs:"notification_id"`
	Destination     string                   `db:"destination" structs:"destination"`
	DestinationType string                   `db:"destination_type" structs:"destination_type"`
	Status          resources.DeliveryStatus `db:"status" structs:"status"`
	SentAt          *time.Time               `db:"sent_at" structs:"sent_at"`
}
