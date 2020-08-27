package horizon

import (
	"encoding/json"
	errors2 "errors"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/distributed_lab/json-api-connector/cerrors"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	settingsTopicsAvailability = "notifications_topics_availability"
	settingsChannels           = "notifications_channels_priority"
	settingsLocale             = "locale"
)

func (c *Connector) IsTopicEnabled(accountId, topic string) (bool, error) {
	item, err := c.GetSettingsItem(accountId, settingsTopicsAvailability)
	if err != nil {
		return false, errors.Wrap(err, "failed to get enabled topics settings")
	}
	// All notifications enabled by default
	if item == nil {
		return true, nil
	}

	var availability map[string]bool
	err = json.Unmarshal(item.Attributes.Value, &availability)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshal notification settings")
	}

	// Check if user disabled all topics
	allEnabled, ok := availability["all_topics"]
	if ok && !allEnabled {
		return false, nil
	}

	enabled, ok := availability[topic]
	if !ok {
		return true, nil
	}

	return enabled, nil
}

func (c *Connector) GetChannels(accountId string) ([]string, error) {
	item, err := c.GetSettingsItem(accountId, settingsChannels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enabled topics settings")
	}
	if item == nil {
		return nil, nil
	}

	var channels []string
	err = json.Unmarshal(item.Attributes.Value, &channels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal notification settings")
	}

	return channels, nil
}

func (c *Connector) GetLocale(accountId string) (*string, error) {
	item, err := c.GetSettingsItem(accountId, settingsLocale)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get enabled topics settings")
	}
	if item == nil {
		return nil, nil
	}

	var locale string
	err = json.Unmarshal(item.Attributes.Value, &locale)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal notification settings")
	}

	return &locale, nil
}

func (c *Connector) GetSettingsItem(accountId, key string) (*SettingsItem, error) {
	path, err := url.Parse(fmt.Sprintf("identities/%s/settings/%s", accountId, key))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse path")
	}

	var settingsItem SettingsItemResponse
	err = c.connector.Get(path, &settingsItem)
	if err != nil {
		var internalErr cerrors.Error
		if errors2.As(err, &internalErr) {
			if internalErr.Status() == http.StatusNotFound {
				return nil, nil
			}
		}
		return nil, errors.Wrap(err, "failed to send request")
	}
	if isJsonNull(settingsItem.Data.Attributes.Value) {
		return nil, nil
	}

	return &settingsItem.Data, nil
}

func isJsonNull(json json.RawMessage) bool {
	return string(json) == "null"
}

type SettingsItemResponse struct {
	Data SettingsItem
}

type SettingsItem struct {
	Attributes SettingsItemAttributes
}

type SettingsItemAttributes struct {
	Key   string
	Value json.RawMessage
}
