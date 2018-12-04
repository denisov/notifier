package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/denisov/notifier/kengusite"
	"github.com/denisov/notifier/shkolanso"
	"github.com/denisov/notifier/telegram"
)

func main() {
	kenguParser := kengusite.NewParser(os.Getenv("KENGU_LOGIN"), os.Getenv("KENGU_PASSWORD"))
	shkolaParser := shkolanso.NewParser(os.Getenv("SHKOLANSO_LOGIN"), os.Getenv("SHKOLANSO_PASSWORD"))

	bot, err := telegram.NewBot(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		fmt.Sprintf("https://%s/webhook", os.Getenv("HOST")),
		kenguParser,
		shkolaParser,
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	http.HandleFunc("/webhook", bot.Handler)
	http.HandleFunc("/cron", authMiddleware(bot.CronHandler, os.Getenv("CRON_KEY")))

	http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)
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
