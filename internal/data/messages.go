package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"
)

const (
	NotificationMessageTemplate = "notification-message-template"
)

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

type TemplateMessageAttributes struct {
	Payload *json.RawMessage `json:"payload"`
	Locale  *string          `json:"locale"`
	Files   []string         `json:"files"`
}
