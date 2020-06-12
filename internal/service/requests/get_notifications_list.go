package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetNotificationsListRequest struct {
	pgdb.OffsetPageParams
	FilterToken              []string `filter:"token"`
	FilterTopic              []string `filter:"topic"`
	FilterDestinationAccount *string  `filter:"destination_account"`
	IncludeDeliveries        bool     `include:"deliveries"`
}

func NewGetNotificationsListRequest(r *http.Request) (GetNotificationsListRequest, error) {
	request := GetNotificationsListRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
