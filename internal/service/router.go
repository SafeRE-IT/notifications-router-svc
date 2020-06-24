package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data/pg"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/horizon"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/handlers"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	horizonConnector := horizon.NewConnector(s.cfg.Client())

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxNotificationsQ(pg.NewNotificationsQ(s.cfg.DB())),
			handlers.CtxDeliveriesQ(pg.NewDeliveriesQ(s.cfg.DB())),
			handlers.CtxHorizon(horizonConnector),
			handlers.CtxDoorman(doorman.New(
				s.cfg.SkipSignCheck(),
				horizonConnector),
			),
			handlers.CtxServices(s.services),
		),
	)

	r.Route("/integrations/notifications", func(r chi.Router) {
		r.Post("/", handlers.CreateNotification)
		r.Get("/", handlers.GetNotificationsList)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetNotification)
			r.Patch("/cancel", handlers.CancelNotification)
		})
		r.Post("/register", handlers.RegisterService)
	})

	return r
}
