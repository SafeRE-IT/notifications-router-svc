package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/helpers"

	"github.com/SafeRE-IT/notifications-router-svc/internal/data"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/router/requests"
)

func GetNotification(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetNotificationRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO: Is allowed for destination
	if !helpers.IsAllowed(r, w) {
		return
	}

	notification, err := helpers.NotificationsQ(r).FilterByID(request.NotificationID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get notification from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if notification == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	deliveries, err := helpers.DeliveriesQ(r).FilterByNotificationID(request.NotificationID).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get deliveries from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	result := resources.NotificationResponse{
		Data: newNotificationModel(*notification, deliveries),
	}
	if request.IncludeDeliveries {
		result.Included = newNotificationIncluded(deliveries)
	}

	ape.Render(w, result)
}

func newNotificationIncluded(deliveries []data.Delivery) resources.Included {
	result := resources.Included{}
	for _, item := range deliveries {
		resource := newDeliveryModel(item)
		result.Add(&resource)
	}
	return result
}

func newDeliveryModel(delivery data.Delivery) resources.Delivery {
	return resources.Delivery{
		Key: resources.NewKeyInt64(delivery.ID, resources.NOTIFICATION_DELIVERY),
		Attributes: resources.DeliveryAttributes{
			Destination:     delivery.Destination,
			DestinationType: delivery.DestinationType,
			SentAt:          delivery.SentAt,
			Status:          delivery.Status,
		},
	}
}
