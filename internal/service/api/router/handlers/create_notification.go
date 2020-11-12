package handlers

import (
	"net/http"
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/helpers"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/router/requests"
)

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateNotificationRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.IsAllowed(r, w) {
		return
	}

	var resultNotification data.Notification
	var resultDeliveries []data.Delivery

	err = helpers.NotificationsQ(r).Transaction(func(q data.NotificationsQ) error {
		notification := data.Notification{
			Topic:   request.Data.Attributes.Topic,
			Token:   request.Data.Attributes.Token,
			Channel: request.Data.Attributes.Channel,
			Message: data.Message(request.Data.Attributes.Message),
		}

		if request.Data.Attributes.ScheduledFor != nil {
			notification.ScheduledFor = *request.Data.Attributes.ScheduledFor
		} else {
			notification.ScheduledFor = time.Now().UTC()
		}

		if request.Data.Attributes.Priority != nil {
			notification.Priority = *request.Data.Attributes.Priority
		} else {
			notification.Priority = resources.NotificationsPriorityMedium
		}

		resultNotification, err = q.Insert(notification)
		if err != nil {
			return errors.Wrap(err, "failed to insert notification")
		}

		deliveries := make([]data.Delivery, len(request.Data.Relationships.Destinations.Data))
		for i, destination := range request.Data.Relationships.Destinations.Data {
			deliveries[i] = data.Delivery{
				NotificationID:  resultNotification.ID,
				Destination:     destination.ID,
				DestinationType: string(destination.Type), // TODO: Use string instead of relation type
				Status:          resources.DeliveryStatusNotSent,
			}
		}

		resultDeliveries, err = q.InsertDeliveries(deliveries)
		if err != nil {
			return errors.Wrap(err, "failed to insert delivery")
		}

		return nil
	})
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create notification")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.NotificationResponse{
		Data: newNotificationModel(resultNotification, resultDeliveries),
	}
	ape.Render(w, result)
}

func newNotificationModel(notification data.Notification, deliveries []data.Delivery) resources.Notification {
	result := resources.Notification{
		Key: resources.NewKeyInt64(notification.ID, resources.NOTIFICATION),
		Attributes: resources.NotificationAttributes{
			CreatedAt:    notification.CreatedAt,
			ScheduledFor: notification.ScheduledFor,
			Topic:        notification.Topic,
			Token:        notification.Token,
			Channel:      notification.Channel,
			Priority:     notification.Priority,
			Message:      resources.Message(notification.Message),
		},
		Relationships: resources.NotificationRelationships{
			Deliveries: resources.RelationCollection{
				Data: make([]resources.Key, len(deliveries)),
			},
		},
	}

	for i, item := range deliveries {
		result.Relationships.Deliveries.Data[i] = resources.NewKeyInt64(item.ID, resources.NOTIFICATION_DELIVERY)
	}

	return result
}
