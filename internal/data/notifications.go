package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"
)

type NotificationsQ interface {
	New() NotificationsQ

	Get() (*Notification, error)
	Select() ([]Notification, error)

	Transaction(fn func(q NotificationsQ) error) error
}

type Notification struct {
	ID           int64     `db:"id" structs:"-"`
	CreatedAt    time.Time `db:"created_at" structs:"created_at"`
	ScheduledFor time.Time `db:"scheduled_for" structs:"scheduled_for"`
	Topic        string    `db:"topic" structs:"topic"`
	Token        *string   `db:"token" structs:"token"`
	Priority     int32     `db:"priority" structs:"priority"`
	Channel      *string   `db:"delivery_type" structs:"delivery_type"`
	Message      Message   `db:"message" structs:"message"`
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
