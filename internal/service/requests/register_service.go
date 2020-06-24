package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type RegisterServiceRequest struct {
	Endpoint string `json:"endpoint"`
	Channel  string `json:"channel"`
}

func NewRegisterServiceRequest(r *http.Request) (RegisterServiceRequest, error) {
	var request RegisterServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *RegisterServiceRequest) validate() error {
	return validation.Errors{
		"endpoint": validation.Validate(&r.Endpoint, validation.Required),
		"channel":  validation.Validate(&r.Channel, validation.Required),
	}.Filter()
}
