## Задать secret

```
now secret add telegram_bot_token xxxxx

now secret add kengu_login xxxxx
now secret add kengu_password xxxxx

now secret add shkolanso_login xxxxx
now secret add shkolanso_password xxxxx

now secret add cron_key xxxxx
```

## Deoloy
    ./deploy.sh


### Удалить старые версии
```
# удалить совсем всё
now rm -y notifier

# удалить то что без алиасов (не-боевые деплойменты)
now rm notifier --safe --yes
```


[![Build Status](https://travis-ci.com/denisov/notifier.svg?branch=master)](https://travis-ci.com/denisov/notifier)