package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type telegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	ChatID      string
}

func sendTelegramNotification(text, proxy string, token string, chatIds []string) ([]telegramResponse, error) {

	responses := []telegramResponse{}

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return responses, fmt.Errorf("не могу распарсить url proxy: %s", err)
	}
	client := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	for _, chatID := range chatIds {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return responses, fmt.Errorf("не могу создать request %s", err)
		}
		query := request.URL.Query()

		query.Add("text", text)
		query.Add("chat_id", chatID)
		request.URL.RawQuery = query.Encode()
		resp, err := client.Do(request)
		if err != nil {
			return responses, fmt.Errorf("не могу отправить запрос %s", err)
		}

		rawResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return responses, fmt.Errorf("не могу прочитать ответ. HTTP код:%d %s", resp.StatusCode, err)
		}

		responseItem := telegramResponse{ChatID: chatID}
		err = json.Unmarshal(rawResponse, &responseItem)
		if err != nil {
			return responses, fmt.Errorf("не могу разобрать ответ telegram api %d %s %s", resp.StatusCode, rawResponse, err)
		}

		responses = append(responses, responseItem)
	}
	return responses, nil
}
