package service

import (
	"context"
	"sync"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/router"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/api/registration"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/types"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/processor"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/config"
)

func runService(service types.Service, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		if err := service.Run(context.Background()); err != nil {
			panic(err)
		}
	}()
}

func Run(cfg config.Config) {
	notificationsStorage := notificators.NewMemoryNotificationsStorage()
	wg := &sync.WaitGroup{}

	processorService := processor.NewProcessor(cfg, notificationsStorage)
	runService(processorService, wg)

	routerApi := router.NewRouterAPI(cfg)
	runService(routerApi, wg)

	registrationApi := registration.NewRegistrationAPI(cfg, notificationsStorage)
	runService(registrationApi, wg)

	wg.Wait()
}
