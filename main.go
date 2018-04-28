package main

import (
	"flag"
	"log"
)

func main() {

	configPath := flag.String("config", "config.yml", "path to config.yml")
	flag.Parse()

	conf, err := getConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	response, err := getContent(conf.Credentials.Login, conf.Credentials.Password)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := getBalance(response)
	if err != nil {
		log.Fatal(err)
	}
	//balance := "asds"

	telegramResponses, err := sendTelegramNotification(
		"Баланс на обеденной карте: "+balance,
		conf.TelegramBot.Proxy,
		conf.TelegramBot.Token,
		conf.TelegramBot.ChatIds,
	)

	if err != nil {
		log.Fatal(err)
	}
	for _, responseItem := range telegramResponses {
		if responseItem.Ok {
			log.Printf("%s => ok\n", responseItem.ChatID)
		} else {
			log.Printf("%s => %d %s\n", responseItem.ChatID, responseItem.ErrorCode, responseItem.Description)
		}
	}
}
