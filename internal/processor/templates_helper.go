package processor

import (
	"encoding/json"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/templates"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/config"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
)

type templatesHelper struct {
	templatesProvider templates.TemplatesProvider
	notificatorCfg    *config.NotificatorConfig
}

func (h *templatesHelper) buildMessage(channel string, delivery data.Delivery, notification data.Notification) (data.Message, error) {
	if notification.Message.Type != data.NotificationMessageTemplate {
		return notification.Message, nil
	}

	var templateAttrs data.TemplateMessageAttributes
	err := json.Unmarshal(notification.Message.Attributes, &templateAttrs)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to get template")
	}

	locale, err := h.getLocale(delivery, notification)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to get locale")
	}

	rawMes, err := h.templatesProvider.GetTemplate(notification.Topic, channel, locale)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to download template")
	}
	if rawMes == nil {
		return data.Message{}, errors.New("template not found")
	}

	if templateAttrs.Payload != nil {
		rawAttrs, err := interpolate(string(rawMes), *templateAttrs.Payload)
		if err != nil {
			return data.Message{}, errors.Wrap(err, "failed to interpolate template")
		}
		rawMes = rawAttrs
	}

	var result data.Message
	err = json.Unmarshal(rawMes, &result)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to marshal template to message")
	}

	return result, nil
}

// TODO: Get locale: 1. Notification model 2. User settings 3. Default for service
// TODO: Use array of locales with priority instead of one locale
func (h *templatesHelper) getLocale(delivery data.Delivery, notification data.Notification) (string, error) {
	return h.notificatorCfg.DefaultLocale, nil
}
