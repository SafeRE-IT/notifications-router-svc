package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"

	"gitlab.com/tokend/connectors/signed"

	"github.com/jmoiron/sqlx/types"

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
	return &processor{
		log:                 config.Log().WithField("runner", serviceName),
		notificationsQ:      pg.NewNotificationsQ(config.DB()),
		deliveriesQ:         pg.NewDeliveriesQ(config.DB()),
		notificatorCfg:      config.NotificatorConfig(),
		notificatorsStorage: notificatorsStorage,
		horizon:             horizon.NewConnector(config.Client()),
	}
}

type processor struct {
	log                 *logan.Entry
	notificationsQ      data.NotificationsQ
	deliveriesQ         data.DeliveriesQ
	notificatorCfg      *config.NotificatorConfig
	notificatorsStorage notificators.NotificatorsStorage
	horizon             *horizon.Connector
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
	message, _ := p.GetMessage(delivery, *notification, channel)

	// TODO: Get identifier

	p.sendNotification(message, delivery.Destination, channel)

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

type MesAttributes struct {
	Payload types.JSONText `json:"payload"`
	Locale  string         `json:"locale"`
}

func (p *processor) GetMessage(delivery data.Delivery, notification data.Notification, channel string) (mes, error) {
	//TODO: Add notification message
	if notification.Message.Type != data.NotificationMessageTemplate {
		//result, _ := json.Marshal(notification.Message)
		return mes{}, nil
	}

	// TODO: Get locale: 1. Notification model 2. User settings 3. Default for service
	// TODO: Add error handling
	locale, _ := p.GetLocale(delivery, notification)

	result, err := p.horizon.GetTemplate(notification.Topic, channel, locale)
	if err != nil {
		return mes{}, errors.Wrap(err, "failed to download template")
	}

	var mes mes
	err = json.Unmarshal(result, &mes)
	if err != nil {
		return mes, errors.Wrap(err, "failed to unmarshal template")
	}

	var mesAttributes MesAttributes
	err = json.Unmarshal(notification.Message.Attributes, &mesAttributes)
	if err != nil {
		return mes, errors.Wrap(err, "failed to unmarshal payload")
	}

	mes.Body, err = interpolate(mes.Body, mesAttributes.Payload)
	if err != nil {
		return mes, errors.Wrap(err, "failed to interpolate message")
	}

	mes.Title, err = interpolate(mes.Title, mesAttributes.Payload)
	if err != nil {
		return mes, errors.Wrap(err, "failed to interpolate message")
	}

	return mes, nil
}

type mes struct {
	Body  string
	Title string
}

func interpolate(tmpl string, payload types.JSONText) (string, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tmpl)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	p := make(map[string]string)
	if err = json.Unmarshal(payload, &p); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal payload")
	}

	var res bytes.Buffer
	if err = t.Execute(&res, p); err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return res.String(), nil
}

func (p *processor) sendNotification(m mes, destination string, channel string) {
	notificatorService, err := p.notificatorsStorage.GetByChannel(channel)
	if err != nil {
		p.log.WithError(err).Error("failed to get notificator service")
	}
	connector := horizon.NewConnector(signed.NewClient(http.DefaultClient, &notificatorService.Endpoint))

	err = connector.SendMessage(horizon.MessageRequest{
		Data: horizon.Message{
			Attributes: horizon.MessageAttributes{
				Owner: destination,
				Title: m.Title,
				Body:  m.Body,
			},
		},
	})
	if err != nil {
		p.log.WithError(err).Error("failed to send notification")
	}
}

func (p *processor) SetDeliveryStatus(id int64, status data.DeliveryStatus) error {
	_, err := p.deliveriesQ.New().
		FilterById(id).
		SetStatus(status).
		Update()
	return err
}
