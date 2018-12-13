## Задать переменные окружения

```
HOST
TELEGRAM_BOT_TOKEN
KENGU_LOGIN
KENGU_PASSWORD
SHKOLANSO_LOGIN
SHKOLANSO_PASSWORD
CRON_KEY
```

## Deploy
```
./deploy.sh
```

## Logs
```
heroku logs -t --app=notify-bot-andr
```

## Scheduler
```
wget -O - "https://ZZZZ.herokuapp.com/cron?key=YYYYY"
```

`curl` часто выдаёт "out of memory", поэтому `wget`


[![Build Status](https://travis-ci.com/denisov/notifier.svg?branch=master)](https://travis-ci.com/denisov/notifier)