FROM golang:1.18-alpine

RUN apk add --no-cache git build-base

WORKDIR /go/src/github.com/SafeRE-IT/notifications-router-svc
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/notifications-router-svc github.com/SafeRE-IT/notifications-router-svc

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/notifications-router-svc /usr/local/bin/notifications-router-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["notifications-router-svc"]
