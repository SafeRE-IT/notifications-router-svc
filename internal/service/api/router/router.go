package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
	"github.com/SafeRE-IT/notifications-router-svc/internal/config"
	"github.com/SafeRE-IT/notifications-router-svc/internal/connectors/horizon"
	"github.com/SafeRE-IT/notifications-router-svc/internal/data/pg"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/helpers"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/router/handlers"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/types"
)

func NewRouterAPI(cfg config.Config) types.Service {
	return &routerAPI{
		cfg: cfg,
	}
}

type routerAPI struct {
	cfg config.Config
}

func (s *routerAPI) Run(ctx context.Context) error {
	router := router(s.cfg)

	if err := s.cfg.Copus().RegisterChi(router); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.cfg.Listener(), router)
}

func router(cfg config.Config) chi.Router {
	r := chi.NewRouter()
	log := cfg.Log().WithFields(map[string]interface{}{
		"service": "notifications-api",
	})
	horizonConnector := horizon.NewConnector(cfg.Client())
	info, err := horizonConnector.Info()
	if err != nil {
		panic(errors.Wrap(err, "failed to get horizon info"))
	}

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			helpers.CtxLog(log),
			helpers.CtxNotificationsQ(pg.NewNotificationsQ(cfg.DB())),
			helpers.CtxDeliveriesQ(pg.NewDeliveriesQ(cfg.DB())),
			helpers.CtxHorizon(horizonConnector),
			helpers.CtxDoorman(doorman.New(
				cfg.SkipSignCheck(),
				horizonConnector),
			),
			helpers.CtxHorizonInfo(info),
		),
	)

	r.Route("/integrations/notifications", func(r chi.Router) {
		r.Post("/", handlers.CreateNotification)
		r.Get("/", handlers.GetNotificationsList)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handlers.GetNotification)
			r.Patch("/cancel", handlers.CancelNotification)
		})
	})

	return r
}
