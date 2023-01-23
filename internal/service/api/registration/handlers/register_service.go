package handlers

import (
	"net/http"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/registration/requests"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/helpers"

	"github.com/SafeRE-IT/notifications-router-svc/internal/notificators"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RegisterService(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRegisterServiceRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	service := notificators.NotificatorService{
		Endpoint: request.Endpoint,
		Channels: request.Channels,
	}
	helpers.NotificatorsStorage(r).Add(service)

	ape.Render(w, newNotificationServiceResponse(service))
}

func newNotificationServiceResponse(service notificators.NotificatorService) resources.NotificatorServiceResponse {
	return resources.NotificatorServiceResponse{
		Data: resources.NotificatorService{
			Key: resources.Key{
				Type: resources.NOTIFICATOR_SERVICE,
			},
			Attributes: resources.NotificatorServiceAttributes{
				Endpoint: service.Endpoint.String(),
				Channels: service.Channels,
			},
		},
	}
}
