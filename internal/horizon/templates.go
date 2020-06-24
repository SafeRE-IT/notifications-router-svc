package horizon

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *Connector) GetTemplate(topic, channel, locale string) ([]byte, error) {
	path, err := url.Parse(fmt.Sprintf("/templates/%s-%s-%s", channel, topic, locale))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	var result json.RawMessage
	c.connector.Get(path, &result)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}
