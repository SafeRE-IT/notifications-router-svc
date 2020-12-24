package processor

import (
	"encoding/json"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/settings"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/providers/templates"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/config"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"
)

type templatesHelper struct {
	templatesProvider templates.TemplatesProvider
	notificatorCfg    *config.NotificatorConfig
	settingsProvider  settings.SettingsProvider
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

	locale, err := h.getLocale(delivery, templateAttrs)
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

	if len(templateAttrs.Files) > 0 {
		result, err = appendFiles(result, templateAttrs.Files)
		if err != nil {
			return data.Message{}, errors.Wrap(err, "failed to append files to message")
		}
	}

	return result, nil
}

// TODO: Use array of locales with priority instead of one locale
func (h *templatesHelper) getLocale(delivery data.Delivery, templateAttrs data.TemplateMessageAttributes) (string, error) {
	if templateAttrs.Locale != nil {
		return *templateAttrs.Locale, nil
	}

	if delivery.DestinationType == data.NotificationDestinationAccount {
		locale, err := h.settingsProvider.GetLocale(delivery.Destination)
		if err != nil {
			return "", errors.Wrap(err, "failed to get locale from settings")
		}
		if locale != nil {
			return *locale, nil
		}
	}

	return h.notificatorCfg.DefaultLocale, nil
}

func appendFiles(destMessage data.Message, files []string) (data.Message, error) {
	var rawAttrs map[string]interface{}
	err := json.Unmarshal(destMessage.Attributes, &rawAttrs)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to unmarshal message attributes")
	}

	rawAttrs["files"] = files
	destMessage.Attributes, err = json.Marshal(rawAttrs)
	if err != nil {
		return data.Message{}, errors.Wrap(err, "failed to marshal message attributes")
	}

	return destMessage, nil
}
