package horizon

import (
	jsonapi "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/tokend/connectors/lazyinfo"
	"gitlab.com/tokend/connectors/signed"
)

type Connector struct {
	connector *jsonapi.Connector

	*lazyinfo.LazyInfoer
}

func NewConnector(client *signed.Client) *Connector {
	return &Connector{
		connector:  jsonapi.NewConnector(client),
		LazyInfoer: lazyinfo.New(client),
	}
}
