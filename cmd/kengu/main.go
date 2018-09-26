package main

import (
	"log"
	"net/http"
	"os"

	"github.com/denisov/kengu/kengusite"
	"github.com/denisov/kengu/telegram"
)

func main() {
	parser := kengusite.NewParser(os.Getenv("KENGU_LOGIN"), os.Getenv("KENGU_PASSWORD"))

	bot, err := telegram.NewBot(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		"https://telegram-bot-andrey-kengu.now.sh/webhookKengu",
		parser,
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	http.HandleFunc("/webhookKengu", bot.Handler)

	http.ListenAndServe("0.0.0.0:8443", nil)
}
