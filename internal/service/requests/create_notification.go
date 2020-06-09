package requests

import (
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"
)

type CreateNotificationRequest struct {
	Data resources.CreateNotification
}

func NewCreateNotificationRequest(r *http.Request) (CreateNotificationRequest, error) {
	var request CreateNotificationRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateNotificationRequest) validate() error {
	return validation.Errors{
		"/data/":               validation.Validate(&r.Data, validation.Required),
		"/data/attributes/":    validation.Validate(&r.Data.Attributes, validation.Required),
		"/data/relationships/": validation.Validate(&r.Data.Attributes, validation.Required),
		"/data/attributes/topic": validation.Validate(&r.Data.Attributes.Topic, validation.Required,
			validation.Length(3, 100)),
		"/data/attributes/token": validation.Validate(&r.Data.Attributes.Token, validation.Length(3, 255)),
		"/data/attributes/send_time": validation.Validate(&r.Data.Attributes.SendTime,
			validation.Min(time.Now().UTC()).Error("should be UTC time in future")),
		"/data/attributes/locale":        nil, // TODO: Check that it is a valid locale string
		"/data/attributes/priority":      nil, // TODO: Check that it is a valid priority
		"/data/attributes/delivery_type": nil, // TODO: Check that it is a valid delivery type
		"/data/attributes/message":       validation.Validate(&r.Data.Attributes.Message, validation.Required),
		// TODO: Check that it is in supported message types
		"/data/attributes/message/type":       validation.Validate(&r.Data.Attributes.Message.Type, validation.Required),
		"/data/attributes/message/attributes": validation.Validate(&r.Data.Attributes.Message.Attributes, validation.Required),
	}.Filter()
}
