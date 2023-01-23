package processor

import (
	"github.com/SafeRE-IT/notifications-router-svc/internal/connectors/notifications"
	"github.com/SafeRE-IT/notifications-router-svc/internal/notificators"
	"github.com/pkg/errors"
)

type notificatorsConnectorProvider struct {
	notificatorsStorage notificators.NotificatorsStorage
}

func (p *notificatorsConnectorProvider) getByChannel(chanel string) (notifications.NotificationsConnector, error) {
	service, err := p.notificatorsStorage.GetByChannel(chanel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get notificators service")
	}

	return notifications.NewRestNotificationsConnector(service.Endpoint), nil
}
