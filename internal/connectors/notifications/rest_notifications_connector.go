package notifications

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/json-api-connector/base"
	"gitlab.com/distributed_lab/json-api-connector/cerrors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/connectors/signed"
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

	postMoment := time.Now()
	status, response, err := c.client.PostJSON(path, body, context.TODO())
	if err != nil {
		if status < 200 || status >= 300 {
			responseBody := "empty"
			if response != nil {
				responseBody = string(response)
			}

			if cerr, ok := err.(cerrors.Error); ok {
				responseBody = string(cerr.Body())
			}

			return errors.Wrap(err, "failed to make request", logan.F{
				"status":        status,
				"response_body": responseBody,
				"time_spent":    time.Since(postMoment).String(),
			})
		}
		return errors.Wrap(err, "failed to make request")
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
