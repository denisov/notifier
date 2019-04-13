FROM golang:1.12-alpine as base

RUN apk --no-cache add git
WORKDIR /go/src/github.com/denisov/notifier
COPY . .

RUN go get -v ./...

RUN apk --update add ca-certificates
# CGO_ENABLED: don’t need to worry about library dependencies
# -ldflags "-s -w" to strip the debugging information
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o notifier github.com/denisov/notifier/cmd/notifier

FROM alpine:3.8
# wget нужен для вызова эндпоинта-крона. Крон есть не во всех бесплатных хостингах, поэтому крон сделал в виде эндпоинта HTTP сервера
RUN apk add --update \
    wget \
    && rm -rf /var/cache/apk/*
COPY --from=base /go/src/github.com/denisov/notifier/notifier /go-telegram-bot-notifier
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

CMD ["/go-telegram-bot-notifier"]
