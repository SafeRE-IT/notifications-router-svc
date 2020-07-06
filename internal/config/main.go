package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/connectors/signed"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	signed.Clienter
	Doorman
	Notificator
	RegistrationAPIer
}

type config struct {
	getter kv.Getter

	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	signed.Clienter
	Doorman
	Notificator
	RegistrationAPIer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:            getter,
		Databaser:         pgdb.NewDatabaser(getter),
		Copuser:           copus.NewCopuser(getter),
		Listenerer:        comfig.NewListenerer(getter),
		Logger:            comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Clienter:          signed.NewClienter(getter),
		Doorman:           NewDoorman(getter),
		Notificator:       NewNotificator(getter),
		RegistrationAPIer: NewRegistrationAPIer(getter),
	}
}
