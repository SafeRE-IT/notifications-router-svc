package handlers

import (
	"net/http"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/registration/requests"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/helpers"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"

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
		Channel:  request.Channel,
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
				Channel:  service.Channel,
			},
		},
	}
}
