package data

import "time"

type DeliveryStatus string

type DeliveriesQ interface {
	New() DeliveriesQ

	Get() (*Delivery, error)
	Select() ([]Delivery, error)

	Transaction(fn func(q DeliveriesQ) error) error
}

const (
	DeliveryStatusNotSent DeliveryStatus = "not_sent"
	DeliveryStatusFailed  DeliveryStatus = "failed"
	DeliveryStatusSent    DeliveryStatus = "sent"
)

type Delivery struct {
	ID              int64          `db:"id" structs:"-"`
	NotificationID  int64          `db:"notification_id" structs:"notification_id"`
	Destination     string         `db:"destination" structs:"destination"`
	DestinationType string         `db:"destination_type" structs:"destination_type"`
	Status          DeliveryStatus `db:"status" structs:"status"`
	SentAt          *time.Time     `db:"sent_at" structs:"sent_at"`
}
