FROM golang:1.21-alpine as base

RUN apk --no-cache add git
WORKDIR /go/src/github.com/denisov/notifier
COPY . .

# CGO_ENABLED: don’t need to worry about library dependencies
# -ldflags "-s -w" to strip the debugging information
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o notifier github.com/denisov/notifier/cmd/notifier

FROM alpine:3.18
COPY --from=base /go/src/github.com/denisov/notifier/notifier /go-telegram-bot-notifier

CMD ["/go-telegram-bot-notifier"]
