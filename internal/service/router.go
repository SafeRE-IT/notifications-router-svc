package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data/pg"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/service/handlers"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxNotificationsQ(pg.NewNotificationsQ(s.cfg.DB())),
			handlers.CtxDeliveriesQ(pg.NewDeliveriesQ(s.cfg.DB())),
		),
	)

	r.Route("/integrations/notifications", func(r chi.Router) {
		r.Post("/", handlers.CreateNotification)
	})

	return r
}
