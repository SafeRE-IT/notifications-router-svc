package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"
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
}

type NotificationPriority int32

const (
	NotificationsPriorityLowest NotificationPriority = iota + 1
	NotificationsPriorityLow
	NotificationsPriorityMedium
	NotificationsPriorityHigh
	NotificationsPriorityHighest
)

const (
	NotificationDestinationAccount = "notification-destination-account"
)

const (
	NotificationMessageTemplate = "notification-message-template"
)

type Notification struct {
	ID           int64                `db:"id" structs:"-"`
	CreatedAt    time.Time            `db:"created_at" structs:"created_at"`
	ScheduledFor time.Time            `db:"scheduled_for" structs:"scheduled_for"`
	Topic        string               `db:"topic" structs:"topic"`
	Token        *string              `db:"token" structs:"token"`
	Priority     NotificationPriority `db:"priority" structs:"priority"`
	Channel      *string              `db:"channel" structs:"channel"`
	Message      Message              `db:"message" structs:"-"`
}

type Message resources.Message

func (m Message) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	return j, err
}

func (m *Message) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, m)
}
