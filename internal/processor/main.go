package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/templates"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/identifier"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/horizon"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/distributed_lab/running"

	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data/pg"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/config"
)

type Processor interface {
	Run(ctx context.Context)
}

const (
	serviceName = "notifications-processor"
)

func NewProcessor(config config.Config, notificatorsStorage notificators.NotificatorsStorage) Processor {
	horizonConnector := horizon.NewConnector(config.Client())
	return &processor{
		log:                 config.Log().WithField("runner", serviceName),
		notificationsQ:      pg.NewNotificationsQ(config.DB()),
		deliveriesQ:         pg.NewDeliveriesQ(config.DB()),
		notificatorCfg:      config.NotificatorConfig(),
		notificatorsStorage: notificatorsStorage,
		horizon:             horizonConnector,
		identifierProvider:  horizonConnector,
		templatesProvider:   templates.NewHorizonTemplatesProvider(config.Client()),
	}
}

type processor struct {
	log                 *logan.Entry
	notificationsQ      data.NotificationsQ
	deliveriesQ         data.DeliveriesQ
	notificatorCfg      *config.NotificatorConfig
	notificatorsStorage notificators.NotificatorsStorage
	horizon             *horizon.Connector
	identifierProvider  identifier.IdentifierProvider
	templatesProvider   templates.TemplatesProvider
}

func (p *processor) Run(ctx context.Context) {
	running.WithBackOff(ctx, p.log,
		serviceName,
		p.processNotifications,
		10*time.Second,
		10*time.Second,
		10*time.Second,
	)
}

func (p *processor) processNotifications(ctx context.Context) error {
	deliveries, err := p.getPendingDeliveries()
	if err != nil {
		return errors.Wrap(err, "failed to get pending deliveries")
	}

	for _, delivery := range deliveries {
		p.log.Info("processing notification", map[string]interface{}{
			"delivery_id":          delivery.ID,
			"notification_id":      delivery.NotificationID,
			"delivery_destination": delivery.Destination,
		})
		err = p.processDelivery(delivery)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to process delivery"),
				map[string]interface{}{
					"delivery_id":          delivery.ID,
					"notification_id":      delivery.NotificationID,
					"delivery_destination": delivery.Destination,
				})
		}
	}

	return nil
}

func (p *processor) getPendingDeliveries() ([]data.Delivery, error) {
	return p.deliveriesQ.New().
		JoinNotification().
		FilterByStatus(data.DeliveryStatusNotSent).
		FilterByScheduledBefore(time.Now().UTC()).
		OrderByPriority(pgdb.OrderTypeDesc).
		Select()
}

func (p *processor) processDelivery(delivery data.Delivery) error {
	// TODO: Join to delivery
	notification, err := p.notificationsQ.New().
		FilterByID(delivery.NotificationID).
		Get()
	if err != nil {
		return errors.Wrap(err, "failed to get notification")
	}
	if notification == nil {
		return errors.New("failed to find notification for delivery")
	}
	p.log.Infof("sending notification %d, notification destination: %d", notification.ID, delivery.Destination)

	// TODO: Check user settings if notification is disabled

	// TODO: Add error handling
	// TODO: Get channel based on available identificator
	channel, _ := p.GetChannel(delivery, *notification)
	message, err := p.GetMessage(delivery, *notification, channel)
	if err != nil {
		return errors.Wrap(err, "failed to create message from template")
	}

	id, err := p.GetIdentifier(channel, delivery)
	if err != nil {
		return errors.Wrap(err, "failed to get identifier")
	}

	println(channel)
	println(string(mustMarshalJSON(message)))
	println(string(mustMarshalJSON(id)))
	//p.sendNotification(message, id.ID, channel)

	if err = p.SetDeliveryStatus(delivery.ID, data.DeliveryStatusSent); err != nil {
		return errors.Wrap(err, "failed to mark delivery sent")
	}

	return nil
}

func (p *processor) GetChannel(delivery data.Delivery, notification data.Notification) (string, error) {
	if notification.Channel != nil {
		return *notification.Channel, nil
	}

	// TODO: Get from user settings

	return p.notificatorCfg.DefaultChannel, nil
}

func (p *processor) GetLocale(delivery data.Delivery, notification data.Notification) (string, error) {
	return p.notificatorCfg.DefaultLocale, nil
}

func (p *processor) GetIdentifier(channel string, delivery data.Delivery) (identifier.Identifier, error) {
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

func (p *processor) GetMessage(delivery data.Delivery, notification data.Notification, channel string) (data.Message, error) {
	if notification.Message.Type != data.NotificationMessageTemplate {
		return notification.Message, nil
	}

	var templateAttrs data.TemplateMessageAttributes
	err := json.Unmarshal(notification.Message.Attributes, &templateAttrs)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to get template")
	}

	// TODO: Get locale: 1. Notification model 2. User settings 3. Default for service
	locale, err := p.GetLocale(delivery, notification)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to get locale")
	}

	rawMes, err := p.templatesProvider.GetTemplate(notification.Topic, channel, locale)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to download template")
	}
	if rawMes == nil {
		return data.Message{}, errors.New("template not found")
	}

	if templateAttrs.Payload != nil {
		rawAttrs, err := interpolate(string(rawMes), *templateAttrs.Payload)
		if err != nil {
			return data.Message{}, errors.Wrap(err, "failed to interpolate template")
		}
		rawMes = rawAttrs
	}

	var result data.Message
	err = json.Unmarshal(rawMes, &result)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to marshal template to message")
	}

	return result, nil
}

func interpolate(tmpl string, payload []byte) ([]byte, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tmpl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	p := make(map[string]interface{})
	if err = json.Unmarshal(payload, &p); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload")
	}

	var res bytes.Buffer
	if err = t.Execute(&res, p); err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return res.Bytes(), nil
}

//func (p *processor) sendNotification(m mes, destination string, channel string) {
//	notificatorService, err := p.notificatorsStorage.GetByChannel(channel)
//	if err != nil {
//		p.log.WithError(err).Error("failed to get notificator service")
//	}
//	connector := horizon.NewConnector(signed.NewClient(http.DefaultClient, &notificatorService.Endpoint))
//
//	err = connector.SendMessage(horizon.MessageRequest{
//		Data: horizon.Message{
//			Attributes: horizon.MessageAttributes{
//				Owner: destination,
//				Title: m.Title,
//				Body:  m.Body,
//			},
//		},
//	})
//	if err != nil {
//		p.log.WithError(err).Error("failed to send notification")
//	}
//}

func (p *processor) SetDeliveryStatus(id int64, status data.DeliveryStatus) error {
	_, err := p.deliveriesQ.New().
		FilterById(id).
		SetStatus(status).
		Update()
	return err
}

func mustMarshalJSON(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return result
}
