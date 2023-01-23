package requests

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/SafeRE-IT/notifications-router-svc/resources"

	"github.com/go-ozzo/ozzo-validation/is"

	"gitlab.com/distributed_lab/logan/v3/errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterServiceRequest struct {
	Endpoint url.URL
	Channels []string
}

func NewRegisterServiceRequest(r *http.Request) (RegisterServiceRequest, error) {
	var request resources.NotificatorServiceResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return RegisterServiceRequest{}, errors.Wrap(err, "failed to unmarshal")
	}

	if err := validateRegisterServiceRequest(request); err != nil {
		return RegisterServiceRequest{}, err
	}

	endpoint, err := url.Parse(request.Data.Attributes.Endpoint)
	if err != nil {
		return RegisterServiceRequest{}, validation.Errors{
			"data/attributes/endpoint": errors.New("must be valid url"),
		}
	}

	return RegisterServiceRequest{
		Endpoint: *endpoint,
		Channels: request.Data.Attributes.Channels,
	}, nil
}

func validateRegisterServiceRequest(r resources.NotificatorServiceResponse) error {
	return validation.Errors{
		"data/attributes/endpoint": validation.Validate(&r.Data.Attributes.Endpoint, validation.Required, is.RequestURI),
		"data/attributes/channels": validation.Validate(&r.Data.Attributes.Channels, validation.Required,
			validation.Length(1, 100)),
	}.Filter()
}
