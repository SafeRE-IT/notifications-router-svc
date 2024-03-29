package horizon

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type MessageRequest struct {
	Data    Message `json:"data"`
	Channel string  `json:"channel"`
}

type Message struct {
	Attributes MessageAttributes `json:"attributes"`
}

type MessageAttributes struct {
	Owner string `json:"owner"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *Connector) SendMessage(m MessageRequest) error {
	path, err := url.Parse("/firebase/notification")
	if err != nil {
		return errors.Wrap(err, "failed to create url")
	}

	var raw json.RawMessage
	err = c.connector.PostJSON(path, m, context.TODO(), &raw)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
