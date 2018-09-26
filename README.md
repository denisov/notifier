## Задать secret

```
now secret add telegram_bot_token xxxxx
now secret add kengu_login xxxxx
now secret add kengu_password xxxxx
```

## Deoloy
    ./deploy.sh


### Удалить старые версии
```
# удалить совсем всё
now rm -y kengu

# удалить то что без алиасов (не-боевые деплойменты)
now rm kengu --safe --yes
```