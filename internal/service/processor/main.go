package processor

import (
	"context"
	"time"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"github.com/SafeRE-IT/notifications-router-svc/internal/providers/settings"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/types"

	"github.com/SafeRE-IT/notifications-router-svc/internal/providers/templates"

	"github.com/SafeRE-IT/notifications-router-svc/internal/providers/identifier"

	"github.com/SafeRE-IT/notifications-router-svc/internal/notificators"

	"github.com/SafeRE-IT/notifications-router-svc/internal/connectors/horizon"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/running"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/SafeRE-IT/notifications-router-svc/internal/data"

	"github.com/SafeRE-IT/notifications-router-svc/internal/config"
)

const (
	serviceName = "notifications-processor"
)

func NewProcessor(config config.Config, notificatorsStorage notificators.NotificatorsStorage) types.Service {
	horizonConnector := horizon.NewConnector(config.Client())
	return &processor{
		log:            config.Log().WithField("runner", serviceName),
		notificatorCfg: config.NotificatorConfig(),
		querier:        newQuerier(config.DB()),
		notificationsConnectorProvider: &notificatorsConnectorProvider{
			notificatorsStorage: notificatorsStorage,
		},
		identifierProvider: horizonConnector,
		templatesHelper: &templatesHelper{
			notificatorCfg:    config.NotificatorConfig(),
			templatesProvider: templates.NewHorizonTemplatesProvider(config.Client()),
			settingsProvider:  horizonConnector,
		},
		settingsProvider: horizonConnector,
	}
}

type processor struct {
	log                            *logan.Entry
	querier                        *querier
	notificatorCfg                 *config.NotificatorConfig
	notificationsConnectorProvider *notificatorsConnectorProvider
	identifierProvider             identifier.IdentifierProvider
	templatesHelper                *templatesHelper
	settingsProvider               settings.SettingsProvider
}

func (p *processor) Run(ctx context.Context) error {
	p.log.Info("started processor")
	running.WithBackOff(ctx, p.log,
		serviceName,
		p.processNotifications,
		3*time.Second,
		3*time.Second,
		3*time.Second,
	)
	return nil
}

func (p *processor) processNotifications(ctx context.Context) error {
	deliveries, err := p.querier.getPendingDeliveries()
	if err != nil {
		return errors.Wrap(err, "failed to get pending deliveries")
	}

	for _, delivery := range deliveries {
		p.log.WithFields(getLoggerFields(delivery)).
			Info("processing notification")
		err = p.processDelivery(delivery)
		if err != nil {
			p.log.WithFields(getLoggerFields(delivery)).WithError(err).
				Error("failed to send to notification, marking it as failed")
			if err := p.querier.setDeliveryStatus(delivery.ID, resources.DeliveryStatusFailed); err != nil {
				return errors.Wrap(err, "failed to set delivery status")
			}
		}
	}

	return nil
}

func (p *processor) processDelivery(delivery data.Delivery) error {
	notification, err := p.querier.getNotification(delivery.NotificationID)
	if err != nil {
		return errors.Wrap(err, "failed to get notification")
	}

	if delivery.DestinationType == data.NotificationDestinationAccount {
		enabled, err := p.settingsProvider.IsTopicEnabled(delivery.Destination, notification.Topic)
		if err != nil {
			return errors.Wrap(err, "failed to check if topic is available")
		}
		if !enabled {
			p.log.WithFields(getLoggerFields(delivery)).
				Info("notification is disabled in user settings, skip sending")
			err = p.querier.setDeliveryStatus(delivery.ID, resources.DeliveryStatusSkipped)
			if err != nil {
				return errors.Wrap(err, "failed to mark delivery skipped")
			}
			return nil
		}
	}

	channelsList, err := p.getChannels(delivery, notification)
	if err != nil {
		return errors.Wrap(err, "failed to get channel")
	}

	for _, channel := range channelsList {
		err = p.sendNotification(channel, delivery, notification)
		if err != nil {
			p.log.WithFields(getLoggerFields(delivery)).
				WithError(err).
				Warnf("failed to send notification with channel - %s, try next channel", channel)
			continue
		}

		if err := p.querier.setDeliveryStatus(delivery.ID, resources.DeliveryStatusSent); err != nil {
			return errors.Wrap(err, "failed to set delivery status")
		}
		return nil
	}

	return errors.New("failed to send notification via all available channels")
}

func (p *processor) sendNotification(channel string, delivery data.Delivery, notification data.Notification) error {
	message, err := p.templatesHelper.buildMessage(channel, delivery, notification)
	if err != nil {
		return errors.Wrap(err, "failed to create message from template")
	}

	id, err := p.getIdentifier(channel, delivery)
	if err != nil {
		return errors.Wrap(err, "failed to get identifier")
	}

	connector, err := p.notificationsConnectorProvider.getByChannel(channel)
	if err != nil {
		return errors.Wrap(err, "failed to get notifications connector")
	}

	p.log.WithFields(map[string]interface{}{
		"channel": channel,
		"message": message,
	}).Debug("Sending message")

	err = connector.SendNotification(id, message, channel)
	if err != nil {
		return errors.Wrap(err, "failed to send notification")
	}

	return nil
}

func (p *processor) getChannels(delivery data.Delivery, notification data.Notification) ([]string, error) {
	if notification.Channel != nil {
		return []string{*notification.Channel}, nil
	}

	if delivery.DestinationType == data.NotificationDestinationAccount {
		channels, err := p.settingsProvider.GetChannels(delivery.Destination)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get channels priority from settings")
		}
		if len(channels) > 0 {
			return channels, nil
		}
	}

	return p.notificatorCfg.DefaultChannelsPriority, nil
}

func (p *processor) getIdentifier(channel string, delivery data.Delivery) (identifier.Identifier, error) {
	if delivery.DestinationType != data.NotificationDestinationAccount {
		return identifier.Identifier{
			ID:   delivery.Destination,
			Type: delivery.DestinationType,
		}, nil
	}
	id, err := p.identifierProvider.GetIdentifierByChannel(channel, delivery.Destination)
	if err != nil {
		return identifier.Identifier{}, err
	}
	if id == nil {
		return identifier.Identifier{
			ID:   delivery.Destination,
			Type: delivery.DestinationType,
		}, nil
	}

	return *id, nil
}

func getLoggerFields(delivery data.Delivery) map[string]interface{} {
	return map[string]interface{}{
		"delivery_id":     delivery.ID,
		"notification_id": delivery.NotificationID,
		"destination":     delivery.Destination,
	}
}
