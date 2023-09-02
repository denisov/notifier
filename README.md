# Уведомления о балансе на обеденной карте


## Github registry
Токен создаётся тут: https://github.com/settings/tokens

Залогиниться в github registry
```
echo $CR_PAT | docker login ghcr.io -u denisov --password-stdin
```
Собрать и запушить image в registry
```
docker build -t ghcr.io/denisov/notifier .
docker push ghcr.io/denisov/notifier
```

Использование
```bash
apt install docker.io

# Задать свои параметры в env-файле докера
echo "HOST=mydomain.ru:8443
KENGU_LOGIN=
KENGU_PASSWORD=
SHKOLANSO_LOGIN=
SHKOLANSO_PASSWORD=
TELEGRAM_BOT_TOKEN=
PORT=8443
CRON_KEY=" > env.txt

# Сгенерить self-signed серт. Обязательно указать свой домен в CN=mydomain.ru (!)
openssl req -newkey rsa:2048 -sha256 -nodes -keyout key.pem -x509 -days 3650 -out cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=mydomain.ru"

# если надо обновить image
docker pull ghcr.io/denisov/notifier


docker restart notifier

# запуск как демон
docker run \
    -d \
    --name notifier \
    --env-file env.txt \
    --mount type=bind,source=/root/cert.pem,target=/cert.pem \
    --mount type=bind,source=/root/key.pem,target=/key.pem \
    -p 8443:8443 \
    ghcr.io/denisov/notifier

# смотреть логи
docker logs -f notifier
```

## Запуск по расписанию
Скопировать `notifier.service` и `notifier.timer` в `/etc/systemd/system`

```bash
systemctl enable notifier.timer
systemctl start notifier.timer

# логи
journalctl -u notifier.timer
journalctl -u notifier.service
```


[![Build Status](https://travis-ci.com/denisov/notifier.svg?branch=master)](https://travis-ci.com/denisov/notifier)