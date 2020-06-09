package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/requests"
)

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewCreateNotificationRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

}
