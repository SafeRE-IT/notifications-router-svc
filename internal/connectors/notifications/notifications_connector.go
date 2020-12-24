package notifications

import (
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/identifier"
)

type NotificationsConnector interface {
	SendNotification(destination identifier.Identifier, message data.Message, channel string) error
}
