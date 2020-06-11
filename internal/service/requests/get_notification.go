package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"gitlab.com/distributed_lab/urlval"
)

type GetNotificationRequest struct {
	NotificationID    int64 `url:"-"`
	IncludeDeliveries bool  `include:"deliveries"`
}

func NewGetNotificationRequest(r *http.Request) (GetNotificationRequest, error) {
	request := GetNotificationRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.NotificationID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
