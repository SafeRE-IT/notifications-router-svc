package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
)

type CancelNotificationRequest struct {
	NotificationID int64 `url:"-"`
}

func NewCancelNotificationRequest(r *http.Request) (CancelNotificationRequest, error) {
	request := CancelNotificationRequest{}

	request.NotificationID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
