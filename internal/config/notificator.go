package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type NotificatorConfig struct {
	DefaultChannelsPriority []string `fig:"default_channels_priority,required"`
	DefaultLocale           string   `fig:"default_locale"`
}

type Notificator interface {
	NotificatorConfig() *NotificatorConfig
}

func NewNotificator(getter kv.Getter) Notificator {
	return &notificator{
		getter: getter,
	}
}

type notificator struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *notificator) NotificatorConfig() *NotificatorConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "notificator")

		config := NotificatorConfig{
			DefaultLocale: "en",
		}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &config
	}).(*NotificatorConfig)
}
