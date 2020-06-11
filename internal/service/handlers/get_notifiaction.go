package handlers

import (
	"net/http"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/requests"
)

func GetNotification(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetNotificationRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	notification, err := NotificationsQ(r).FilterByID(request.NotificationID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get notification from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if notification == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	deliveries, err := DeliveriesQ(r).FilterByNotificationID(request.NotificationID).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get deliveries from DB")
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
			Status:          string(delivery.Status), // TODO: Use enums from resources
		},
	}
}
