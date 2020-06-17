package handlers

import (
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/requests"
)

func CancelNotification(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCancelNotificationRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	notification, err := NotificationsQ(r).
		FilterByID(request.NotificationID).
		Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get notification")
		ape.Render(w, problems.InternalError())
		return
	}
	if notification == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	if notification.ScheduledFor.Before(time.Now().UTC()) {
		ape.Render(w, problems.BadRequest(validation.Errors{
			"id": errors.New("only notifications scheduled for future can be canceled")}))
		return
	}

	deliveries, err := DeliveriesQ(r).
		SetStatus(data.DeliveryStatusCanceled).
		FilterByNotificationID(notification.ID).
		Update()
	if err != nil {
		Log(r).WithError(err).Error("failed to update deliveries status")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.NotificationResponse{
		Data:     newNotificationModel(*notification, deliveries),
		Included: newNotificationIncluded(deliveries),
	}
	ape.Render(w, response)
}
