package processor

import (
	"github.com/pkg/errors"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/connectors/notifications"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"
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
