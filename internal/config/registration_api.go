package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type RegistrationAPIConfig struct {
	Addr string `fig:"addr,required"`
}

type RegistrationAPIer interface {
	RegistrationAPIConfig() *RegistrationAPIConfig
}

func NewRegistrationAPIer(getter kv.Getter) RegistrationAPIer {
	return &registrationAPIer{
		getter: getter,
	}
}

type registrationAPIer struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *registrationAPIer) RegistrationAPIConfig() *RegistrationAPIConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "registration_api")

		config := RegistrationAPIConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &config
	}).(*RegistrationAPIConfig)
}
