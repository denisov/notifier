set -x

go test ./... && \
    heroku container:push web --app=notify-bot-andr && \
    heroku container:release web --app=notify-bot-andr
