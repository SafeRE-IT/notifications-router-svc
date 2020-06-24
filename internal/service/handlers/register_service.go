package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/requests"
)

func RegisterService(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRegisterServiceRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w) {
		return
	}

	Services(r)[request.Channel] = request.Endpoint

	ape.Render(w, http.StatusNoContent)
}
