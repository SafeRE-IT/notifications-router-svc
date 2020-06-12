package handlers

import (
	"net/http"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/requests"
)

func GetNotificationsList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetNotificationsListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO: Is allowed for destination
	if !isAllowed(r, w) {
		return
	}

	notificationsQ := NotificationsQ(r)
	applyFilters(notificationsQ, request)
	notifications, err := notificationsQ.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get notifications")
		ape.Render(w, problems.InternalError())
		return
	}

	// TODO: Join instead of two selects
	notificationIds := make([]int64, len(notifications))
	for i, item := range notifications {
		notificationIds[i] = item.ID
	}

	deliveriesQ := DeliveriesQ(r).FilterByNotificationID(notificationIds...)
	if request.FilterDestinationAccount != nil {
		deliveriesQ.FilterByDestination(*request.FilterDestinationAccount).
			FilterByDestinationType(data.NotificationDestinationAccount)
	}
	deliveries, err := deliveriesQ.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get deliveries")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.NotificationListResponse{
		Data: newNotificationsList(notifications, deliveries),
	}
	if request.IncludeDeliveries {
		response.Included = newNotificationIncluded(deliveries)
	}
	ape.Render(w, response)
}

func applyFilters(q data.NotificationsQ, request requests.GetNotificationsListRequest) {
	q.Page(request.OffsetPageParams)

	if request.FilterToken != nil {
		q.FilterByToken(request.FilterToken...)
	}

	if request.FilterTopic != nil {
		q.FilterByTopic(request.FilterTopic...)
	}

	if request.FilterDestinationAccount != nil {
		q.FilterByDestination(*request.FilterDestinationAccount, data.NotificationDestinationAccount)
	}
}

func newNotificationsList(notifications []data.Notification, deliveries []data.Delivery) []resources.Notification {
	result := make([]resources.Notification, len(notifications))
	for i, notification := range notifications {
		notificationDeliveries := make([]data.Delivery, 0, len(deliveries))
		for _, delivery := range deliveries {
			if delivery.NotificationID == notification.ID {
				notificationDeliveries = append(notificationDeliveries, delivery)
			}
		}
		result[i] = newNotificationModel(notification, notificationDeliveries)
	}
	return result
}