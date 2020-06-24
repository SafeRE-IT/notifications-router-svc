package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func NewProcessor(config config.Config, services map[string]string) Processor {
	return &processor{
		log:            config.Log().WithField("runner", serviceName),
		notificationsQ: pg.NewNotificationsQ(config.DB()),
		deliveriesQ:    pg.NewDeliveriesQ(config.DB()),
		notificatorCfg: config.NotificatorConfig(),
		services:       services,
	}
}

type processor struct {
	log            *logan.Entry
	notificationsQ data.NotificationsQ
	deliveriesQ    data.DeliveriesQ
	notificatorCfg *config.NotificatorConfig
	services       map[string]string
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
		err = p.processDelivery(delivery)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to process delivery"),
				map[string]interface{}{
					"delivery_id":     delivery.ID,
					"notification_id": delivery.NotificationID,
				})
		}
	}

	return nil
}

func (p *processor) getPendingDeliveries() ([]data.Delivery, error) {
	return p.deliveriesQ.New().
		JoinNotification().
		FilterByStatus(data.DeliveryStatusNotSent).
		FilterByScheduledForAfter(time.Now()).
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

	// TODO: Check user settings if notification is disabled

	channel, err := p.GetChannel(delivery, *notification)

	if notification.Message.Type == data.NotificationMessageTemplate {
		// TODO: Get locale: 1. Notification model 2. User settings 3. Default for service

		// TODO: Get template
	}

	// TODO: Get identifier

	// TODO: Send notification
	rawNotification, _ := json.Marshal(notification)
	p.log.Info(string(rawNotification))
	p.log.Info(channel)
	p.log.Info(p.services[channel])

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

func (p *processor) SetDeliveryStatus(id int64, status data.DeliveryStatus) error {
	_, err := p.deliveriesQ.New().
		FilterById(id).
		SetStatus(status).
		Update()
	return err
}
