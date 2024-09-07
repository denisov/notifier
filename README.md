# Уведомления о балансе на обеденной карте

Обновить либы
```
go get -u ./...
```

Удалить старое из go.sum и go.mod
```
go mod tidy
```


Собрать бинарник
```
~/code/go_prj/notifier ❯ go build github.com/denisov/notifier/cmd/notifier
```



## Github registry
Токен создаётся тут: https://github.com/settings/tokens

Залогиниться в github registry
```
echo $CR_PAT | docker login ghcr.io -u denisov --password-stdin
```
Собрать и запушить image в registry
```
docker buildx build --platform linux/amd64 -t ghcr.io/denisov/notifier .
docker push ghcr.io/denisov/notifier
```

## Использование

Установить docker и compose https://docs.docker.com/engine/install/ubuntu/

Сгенерить self-signed серт. Обязательно указать __свой__ домен в CN=mydomain.ru (!)
```bash
openssl req -newkey rsa:2048 -sha256 -nodes -keyout key.pem -x509 -days 3650 -out cert.pem -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=mydomain.ru"
```

загрузить/обновить image
```
docker pull ghcr.io/denisov/notifier
```

Создать директорию `notifier`, а в ней `compose.yml`

```bash
mkdir -p notifier && cd notifier
echo "
services:
  notifier:
    image: ghcr.io/denisov/notifier
    container_name: notifier
    environment:
      - HOST=********:8443
      - KENGU_LOGIN=***
      - KENGU_PASSWORD=****
      - SHKOLANSO_LOGIN=*****
      - SHKOLANSO_PASSWORD=******
      - TELEGRAM_BOT_TOKEN=*****
      - PORT=8443
      - CRON_KEY=*************
    volumes:
      - /root/cert.pem:/cert.pem
      - /root/key.pem:/key.pem
    ports:
      - 8443:8443
" > compose.yml
```
Установить свои значения в переменных окружения

Запустить как демона (из директории notifier)
```bash
docker compose up -d
```

Остановить (вместе с удалением контейнеров)
```bash
docker compose down
```

Обновить после обновления образа
```bash
docker compose pull
docker compose up --force-recreate --build -d
docker image prune -f
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