package handlers

import (
	"net/http"
	"time"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/helpers"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/router/requests"
)

func CancelNotification(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCancelNotificationRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !helpers.IsAllowed(r, w) {
		return
	}

	notification, err := helpers.NotificationsQ(r).
		FilterByID(request.NotificationID).
		Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get notification")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if notification == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if notification.ScheduledFor.Before(time.Now().UTC()) {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"id": errors.New("only notifications scheduled for future can be canceled")})...)
		return
	}

	deliveries, err := helpers.DeliveriesQ(r).
		SetStatus(resources.DeliveryStatusCanceled).
		FilterByNotificationID(notification.ID).
		Update()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update deliveries status")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := resources.NotificationResponse{
		Data:     newNotificationModel(*notification, deliveries),
		Included: newNotificationIncluded(deliveries),
	}
	ape.Render(w, response)
}
