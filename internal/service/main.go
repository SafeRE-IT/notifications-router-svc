package service

import (
	"context"
	"sync"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/router"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/registration"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/types"

	"github.com/SafeRE-IT/notifications-router-svc/internal/notificators"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/processor"

	"github.com/SafeRE-IT/notifications-router-svc/internal/config"
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
