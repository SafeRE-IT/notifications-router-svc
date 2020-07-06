package notificators

import "net/url"

type NotificatorService struct {
	Endpoint url.URL
	Channels []string
}
