FROM golang:1.12

WORKDIR /go/src/gitlab.com/tokend/notifications/notifications-router-svc
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/notifications-router-svc gitlab.com/tokend/notifications/notifications-router-svc

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/notifications-router-svc /usr/local/bin/notifications-router-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["notifications-router-svc"]
