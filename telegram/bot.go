package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/denisov/kengu"
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	BotAPI   *tgbotapi.BotAPI
	SiteData kengu.DataSource
}

// NewBot создаёт нового бота
func NewBot(token string, webhookURL string, source kengu.DataSource) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "Can't create bot")
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// todo в конфиг
	bot.Debug = true

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		return nil, errors.Wrap(err, "can't set webhook")
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, errors.Wrap(err, "can't get webhook info")
	}
	if info.LastErrorDate != 0 {
		return nil, errors.Wrap(err, "Telegram callback failed")
	}

	return &Bot{
		BotAPI:   bot,
		SiteData: source,
	}, nil
}

// Handler это обработчик webhook запросов
func (bot *Bot) Handler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read:%v", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Printf("Failed to unmarshal:%v", err)
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}

	var response string

	// TODO проверка на ошибки, если текста нет?
	command := update.Message.Text
	// TODO
	if command == "d" {
		response = bot.getDnevnikResponse()
	} else {
		response = bot.getBalanceResponse()
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)

	_, err = bot.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %s", err)
	}
}

func (bot *Bot) getBalanceResponse() string {
	data, err := bot.SiteData.GetData()

	if err != nil {
		return fmt.Sprintf("ERROR: %+v", err)
	}

	return "Баланс на обеденной карте:" + data
}

func (bot *Bot) getDnevnikResponse() string {
	return "Оценки..."
}
