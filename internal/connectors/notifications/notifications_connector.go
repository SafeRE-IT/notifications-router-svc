package notifications

import (
	"github.com/SafeRE-IT/notifications-router-svc/internal/data"
	"github.com/SafeRE-IT/notifications-router-svc/internal/providers/identifier"
)

type NotificationsConnector interface {
	SendNotification(destination identifier.Identifier, message data.Message, channel string) error
}
