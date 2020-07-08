package data

import (
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type NotificationsQ interface {
	New() NotificationsQ

	Get() (*Notification, error)
	Select() ([]Notification, error)

	Transaction(fn func(q NotificationsQ) error) error

	Insert(data Notification) (Notification, error)
	InsertDeliveries(data []Delivery) ([]Delivery, error)

	Page(pageParams pgdb.OffsetPageParams) NotificationsQ

	FilterByID(id ...int64) NotificationsQ
	FilterByToken(tokens ...string) NotificationsQ
	// TODO: Two separate methods
	FilterByDestination(destination string, destinationType string) NotificationsQ
	FilterByTopic(topics ...string) NotificationsQ
	FilterByScheduledAfter(time time.Time) NotificationsQ
	FilterByScheduledBefore(time time.Time) NotificationsQ
}

const (
	NotificationDestinationAccount = "notification-destination-account"
	NotificationDestinationEmail   = "notidication-destination-email"
	NotificationDestinationPhone   = "notification-destination-phone"
)

type Notification struct {
	ID           int64                          `db:"id" structs:"-"`
	CreatedAt    time.Time                      `db:"created_at" structs:"created_at"`
	ScheduledFor time.Time                      `db:"scheduled_for" structs:"scheduled_for"`
	Topic        string                         `db:"topic" structs:"topic"`
	Token        *string                        `db:"token" structs:"token"`
	Priority     resources.NotificationPriority `db:"priority" structs:"priority"`
	Channel      *string                        `db:"channel" structs:"channel"`
	Message      Message                        `db:"message" structs:"-"`
}
