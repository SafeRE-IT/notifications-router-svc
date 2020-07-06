package helpers

import (
	"context"
	"net/http"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/notificators"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/connectors/horizon"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	notificationsQCtxKey
	deliveriesQCtxKey
	horizonCtxKey
	doormanCtxKey
	notificatorsStorageCtxKey
	horizonInfoCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxNotificationsQ(entry data.NotificationsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificationsQCtxKey, entry)
	}
}

func NotificationsQ(r *http.Request) data.NotificationsQ {
	return r.Context().Value(notificationsQCtxKey).(data.NotificationsQ).New()
}

func CtxDeliveriesQ(entry data.DeliveriesQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, deliveriesQCtxKey, entry)
	}
}

func DeliveriesQ(r *http.Request) data.DeliveriesQ {
	return r.Context().Value(deliveriesQCtxKey).(data.DeliveriesQ).New()
}

func CtxHorizon(h *horizon.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonCtxKey, h)
	}
}

func Horizon(r *http.Request) *horizon.Connector {
	return r.Context().Value(horizonCtxKey).(*horizon.Connector)
}

func CtxDoorman(d doorman.Doorman) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, doormanCtxKey, d)
	}
}

func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(doormanCtxKey).(doorman.Doorman)
	return d.Check(r, constraints...)
}

func CtxNotificatorsStorage(v notificators.NotificatorsStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, notificatorsStorageCtxKey, v)
	}
}

func NotificatorsStorage(r *http.Request) notificators.NotificatorsStorage {
	return r.Context().Value(notificatorsStorageCtxKey).(notificators.NotificatorsStorage)
}

func CtxHorizonInfo(v regources.HorizonState) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonInfoCtxKey, v)
	}
}

func HorizonInfo(r *http.Request) regources.HorizonState {
	return r.Context().Value(horizonInfoCtxKey).(regources.HorizonState)
}
