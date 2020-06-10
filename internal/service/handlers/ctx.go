package handlers

import (
	"context"
	"net/http"

	"gitlab.com/tokend/notifications/notifications-router-svc/internal/data"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	notificationsQCtxKey
	deliveriesQCtxKey
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
