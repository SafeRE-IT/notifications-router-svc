package requests

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-ozzo/ozzo-validation/is"

	"gitlab.com/distributed_lab/logan/v3/errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterServiceRequest struct {
	Endpoint    url.URL `json:"-"`
	Channel     string  `json:"channel"`
	RawEndpoint string  `json:"endpoint"`
}

func NewRegisterServiceRequest(r *http.Request) (RegisterServiceRequest, error) {
	var request RegisterServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	if err := request.validate(); err != nil {
		return request, err
	}

	endpoint, err := url.Parse(request.RawEndpoint)
	if err != nil {
		return request, validation.Errors{
			"endpoint": errors.New("must be valid url"),
		}
	}
	request.Endpoint = *endpoint

	return request, nil
}

func (r *RegisterServiceRequest) validate() error {
	return validation.Errors{
		"endpoint": validation.Validate(&r.RawEndpoint, validation.Required, is.RequestURI),
		"channel":  validation.Validate(&r.Channel, validation.Required),
	}.Filter()
}
