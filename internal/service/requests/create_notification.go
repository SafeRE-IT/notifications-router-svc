package requests

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

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
		"/data/":            validation.Validate(&r.Data, validation.Required),
		"/data/attributes/": validation.Validate(&r.Data.Attributes, validation.Required),
		// TODO: Validate destinations
		"/data/relationships/": validation.Validate(&r.Data.Relationships, validation.Required),
		// TODO: Get limit of destinations from config
		"/data/relationships/destinations/data": validation.Validate(&r.Data.Relationships.Destinations.Data,
			validation.Required, validation.Length(1, 100)),
		"/data/attributes/topic": validation.Validate(&r.Data.Attributes.Topic, validation.Required,
			validation.Length(3, 100)),
		"/data/attributes/token": validation.Validate(&r.Data.Attributes.Token, validation.Length(3, 255)),
		"/data/attributes/scheduled_for": validation.Validate(&r.Data.Attributes.ScheduledFor,
			validation.Min(time.Now().UTC()).Error("should be UTC time in future")),
		"/data/attributes/priority": validation.Validate(&r.Data.Attributes.Priority,
			validation.Min(data.NotificationsPriorityLowest),
			validation.Max(data.NotificationsPriorityHighest),
		),
		"/data/attributes/channel": nil, // TODO: Check that it is a valid delivery type
		// TODO: Validate message
		"/data/attributes/message": validation.Validate(&r.Data.Attributes.Message, validation.Required),
		// TODO: Check that it is in supported message types
		"/data/attributes/message/type":       validation.Validate(&r.Data.Attributes.Message.Type, validation.Required),
		"/data/attributes/message/attributes": validation.Validate(&r.Data.Attributes.Message.Attributes, validation.Required),
	}.Filter()
}
