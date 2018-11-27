package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/denisov/notifier/kengusite"
	"github.com/denisov/notifier/telegram"
)

func main() {
	parser := kengusite.NewParser(os.Getenv("KENGU_LOGIN"), os.Getenv("KENGU_PASSWORD"))

	bot, err := telegram.NewBot(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		"https://telegram-bot-andrey-notifier.now.sh/webhookKengu",
		parser,
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	http.HandleFunc("/webhookKengu", bot.Handler)
	http.HandleFunc("/cron", authMiddleware(bot.CronHandler, os.Getenv("CRON_KEY")))

	http.ListenAndServe("0.0.0.0:8443", nil)
}

func authMiddleware(next http.HandlerFunc, authKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key != authKey {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `{"error":"invalid_key"}`)
			return
		}
		next.ServeHTTP(w, r)
	}
}
