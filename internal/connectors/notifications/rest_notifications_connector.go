package notifications

import (
	"context"
	"net/http"
	"net/url"

	"gitlab.com/tokend/connectors/signed"

	"github.com/pkg/errors"

	"gitlab.com/distributed_lab/json-api-connector/base"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/identifier"
)

const createNotificationType = "create-notification"

func NewRestNotificationsConnector(endpoint url.URL) NotificationsConnector {
	return &restNotificationsConnector{
		endpoint: endpoint,
		client:   base.NewConnector(signed.NewClient(http.DefaultClient, &endpoint)),
	}
}

type restNotificationsConnector struct {
	endpoint url.URL
	client   *base.Connector
}

func (c *restNotificationsConnector) SendNotification(identifier identifier.Identifier, message data.Message, channel string) error {
	path, err := url.Parse("notifications")
	if err != nil {
		return errors.Wrap(err, "failed to parse url")
	}

	body := createNotificationBody{
		Data: createNotificationData{
			Type: createNotificationType,
			Attributes: createNotificationAttributes{
				Message: message,
				Channel: channel,
			},
			Relationships: createNotificationRelationships{
				Destination: createNotificationDestination{
					Data: key{
						ID:   identifier.ID,
						Type: identifier.Type,
					},
				},
			},
		},
	}

	status, response, err := c.client.PostJSON(path, body, context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	if status < 200 || status >= 300 {
		responseBody := "empty"
		if response != nil {
			responseBody = string(response)
		}
		return errors.Errorf("request failed, status code - %d, response body - %s", status, responseBody)
	}
	return nil
}

type createNotificationBody struct {
	Data createNotificationData
}

type createNotificationData struct {
	Type          string
	Attributes    createNotificationAttributes
	Relationships createNotificationRelationships
}

type createNotificationAttributes struct {
	Message data.Message
	Channel string `json:"channel"`
}

type createNotificationRelationships struct {
	Destination createNotificationDestination
}

type createNotificationDestination struct {
	Data key
}

type key struct {
	ID   string
	Type string
}
