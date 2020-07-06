package requests

import (
	"encoding/json"
	"net/http"
	"net/url"

	"gitlab.com/tokend/notifications/notifications-router-svc/resources"

	"github.com/go-ozzo/ozzo-validation/is"

	"gitlab.com/distributed_lab/logan/v3/errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterServiceRequest struct {
	Endpoint url.URL
	Channel  string
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
		Channel:  request.Data.Attributes.Channel,
	}, nil
}

func validateRegisterServiceRequest(r resources.NotificatorServiceResponse) error {
	return validation.Errors{
		"data/attributes/endpoint": validation.Validate(&r.Data.Attributes.Endpoint, validation.Required, is.RequestURI),
		"data/attributes/channel":  validation.Validate(&r.Data.Attributes.Channel, validation.Required),
	}.Filter()
}
