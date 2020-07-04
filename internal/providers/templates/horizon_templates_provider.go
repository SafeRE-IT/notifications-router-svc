package templates

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"gitlab.com/distributed_lab/json-api-connector/base"

	"gitlab.com/distributed_lab/json-api-connector/client"
)

func NewHorizonTemplatesProvider(client client.Client) TemplatesProvider {
	return &horizonTemplatesProvider{
		client: base.NewConnector(client),
	}
}

type horizonTemplatesProvider struct {
	client *base.Connector
}

func (c *horizonTemplatesProvider) GetTemplate(topic, channel, locale string) ([]byte, error) {
	path, err := url.Parse(fmt.Sprintf("/templates/%s-%s-%s", channel, topic, locale))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	response, err := c.client.Get(path)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	return response, nil
}
