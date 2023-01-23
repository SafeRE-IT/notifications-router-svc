package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/SafeRE-IT/notifications-router-svc/internal/types"

	"github.com/SafeRE-IT/notifications-router-svc/internal/data"

	validation "github.com/go-ozzo/ozzo-validation"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/SafeRE-IT/notifications-router-svc/resources"
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
	return mergeErrors(validation.Errors{
		"/data/attributes/topic": validation.Validate(&r.Data.Attributes.Topic, validation.Required,
			validation.Length(3, 100)),
		"/data/attributes/token": validation.Validate(&r.Data.Attributes.Token, validation.Length(3, 255)),
		"/data/attributes/scheduled_for": validation.Validate(&r.Data.Attributes.ScheduledFor,
			validation.Min(time.Now().UTC()).Error("should be UTC time in future")),
		"/data/attributes/priority": validation.Validate(&r.Data.Attributes.Priority,
			validation.Min(resources.NotificationsPriorityLowest),
			validation.Max(resources.NotificationsPriorityHighest),
		),
		"/data/attributes/message/type":       validation.Validate(&r.Data.Attributes.Message.Type, validation.Required),
		"/data/attributes/message/attributes": validation.Validate(&r.Data.Attributes.Message.Attributes, validation.Required),
	},
		validateDestinationsList(r.Data.Relationships.Destinations.Data),
	).Filter()
}

func validateDestinationsList(destinations []resources.Key) validation.Errors {
	validationErrors := validation.Errors{
		"/data/relationships/destinations/data": validation.Validate(&destinations,
			validation.Required, validation.Length(1, 100)),
	}

	for i, destination := range destinations {
		validationErrors[fmt.Sprintf("/data/relationships/destinations/data/%d", i)] =
		// TODO: Use string instead of type
			validateDestination(string(destination.Type), destination.ID)
	}

	return validationErrors
}

func validateDestination(destinationType string, destination string) error {
	switch destinationType {
	case data.NotificationDestinationAccount:
		return types.AccountID(destination).Validate()
	case data.NotificationDestinationEmail:
		if !govalidator.IsEmail(destination) {
			return errors.New("must be valid email")
		}
		return nil
	case data.NotificationDestinationPhone:
		return types.ValidatePhoneNumber(destination)
	default:
		return nil
	}
}

func mergeErrors(validationErrors ...validation.Errors) validation.Errors {
	result := make(validation.Errors)
	for _, errs := range validationErrors {
		for key, err := range errs {
			result[key] = err
		}
	}
	return result
}
