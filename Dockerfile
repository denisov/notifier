FROM golang:1.11-alpine as base

RUN apk --no-cache add git
WORKDIR /go/src/github.com/denisov/kengu
COPY . .

RUN go get -v ./...

RUN apk --update add ca-certificates
# CGO_ENABLED: don’t need to worry about library dependencies
# -ldflags "-s -w" to strip the debugging information
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o kengu github.com/denisov/kengu/cmd/kengu 

FROM scratch
COPY --from=base /go/src/github.com/denisov/kengu/kengu /go-telegram-bot-kengu
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8443
CMD ["/go-telegram-bot-kengu"]
