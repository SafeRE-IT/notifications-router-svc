package registration

import (
	"context"
	"net"
	"net/http"

	"github.com/pkg/errors"

	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/registration/handlers"

	"gitlab.com/distributed_lab/ape"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/api/helpers"

	"github.com/go-chi/chi"

	"github.com/SafeRE-IT/notifications-router-svc/internal/config"
	"github.com/SafeRE-IT/notifications-router-svc/internal/notificators"
	"github.com/SafeRE-IT/notifications-router-svc/internal/service/types"
)

func NewRegistrationAPI(cfg config.Config, notificatorsStorage notificators.NotificatorsStorage) types.Service {
	return &registrationAPI{
		cfg:                 cfg,
		notificatorsStorage: notificatorsStorage,
	}
}

type registrationAPI struct {
	cfg                 config.Config
	notificatorsStorage notificators.NotificatorsStorage
}

func (s *registrationAPI) Run(ctx context.Context) error {
	r := s.router()

	listener, err := net.Listen("tcp", s.cfg.RegistrationAPIConfig().Addr)
	if err != nil {
		return errors.Wrap(err, "failed to create listener")
	}
	return http.Serve(listener, r)
}

func (s *registrationAPI) router() chi.Router {
	r := chi.NewRouter()
	log := s.cfg.Log().WithFields(map[string]interface{}{
		"service": "registraation-api",
	})

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			helpers.CtxLog(log),
			helpers.CtxNotificatorsStorage(s.notificatorsStorage),
		),
	)

	r.Route("/services", func(r chi.Router) {
		r.Post("/", handlers.RegisterService)
	})

	return r
}
