package horizon

import (
	"fmt"
	"net/url"

	"github.com/SafeRE-IT/notifications-router-svc/internal/data"

	"github.com/pkg/errors"

	"github.com/SafeRE-IT/notifications-router-svc/internal/providers/identifier"
)

const (
	channelEmail = "email"
	channelPush  = "push"
	channelSms   = "sms"
)

func (c *Connector) GetIdentifierByChannel(channel string, accountId string) (*identifier.Identifier, error) {
	identity, err := c.GetIdentity(accountId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get identity")
	}
	if identity == nil {
		return nil, nil
	}

	switch channel {
	case channelEmail:
		return &identifier.Identifier{
			Type: data.NotificationDestinationEmail,
			ID:   identity.Attributes.Email,
		}, nil
	case channelSms:
		if identity.Attributes.Phone == nil {
			return nil, nil
		}
		return &identifier.Identifier{
			Type: data.NotificationDestinationPhone,
			ID:   *identity.Attributes.Phone,
		}, nil
	default:
		return nil, nil
	}
}

func (c *Connector) GetIdentity(accountId string) (*IdentityData, error) {
	path, err := url.Parse(fmt.Sprintf("/identities?filter[address]=%s", accountId))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse path")
	}
	var result IdentitiesResponse
	err = c.connector.Get(path, &result)

	if len(result.Data) == 0 {
		return nil, nil
	}

	return &result.Data[0], nil
}

type IdentitiesResponse struct {
	Data []IdentityData `json:"data"`
}

type IdentityData struct {
	Type       string             `json:"type"`
	ID         string             `json:"id"`
	Attributes IdentityAttributes `json:"attributes"`
}

type IdentityAttributes struct {
	Address string  `json:"address"`
	Email   string  `json:"email"`
	Phone   *string `json:"phone_number"`
}
