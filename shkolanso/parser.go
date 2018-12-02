package shkolanso

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

// Parser это парсер сайта школа нсо
type Parser struct {
	login    string
	password string
}

// NewParser создаёт парсер сайта школа нсо
func NewParser(login, password string) Parser {
	return Parser{
		login:    login,
		password: password,
	}
}

const formURL = "https://shkola.nso.ru/auth/login"
const mainURL = "https://shkola.nso.ru"
const summaryMarksURL = "https://shkola.nso.ru/api/MarkService/GetSummaryMarks"

// GetData возвращает данные для уведомления
func (p Parser) GetData() (string, error) {
	summaryMarksRaw, err := getSummaryMarksRaw(p.login, p.password)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer summaryMarksRaw.Close()
	summaryMarksRawContent, err := ioutil.ReadAll(summaryMarksRaw)
	if err != nil {
		return "", err
	}

	// распарсить json ответ в структуру
	sm, err := parseJSONResponse(summaryMarksRawContent)
	if err != nil {
		return "", errors.Wrap(err, "Не могу распарсить ответ")
	}

	// на основе структуры cформировать ответ
	textMessage, err := composeTextMessage(sm, time.Now())
	if err != nil {
		return "", errors.Wrap(err, "Не могу составить ответ")
	}

	return textMessage, nil
}

// parseJSONResponse парсит json ответ
// метод выделен отдельно для удобства тестирования. Тут можно много багов словить при изменении формата ответа
func parseJSONResponse(response []byte) (summaryMarks, error) {
	var sm summaryMarks
	err := json.Unmarshal(response, &sm)
	if err != nil {
		return sm, errors.Wrap(err, "Не могу распарсить ответ")
	}
	return sm, nil
}

// composeTextMessage составляет текстовое сообщение из данных структуры summaryMarks
func composeTextMessage(sm summaryMarks, dt time.Time) (string, error) {
	var message string

	for _, disciplineMarkItem := range sm.DisciplineMarks {
		discipline := disciplineMarkItem.Discipline
		for _, markItem := range disciplineMarkItem.Marks {
			if markItem.Date.Truncate(24 * time.Hour).Equal(dt.Truncate(24 * time.Hour)) {
				if len(message) != 0 {
					message = message + "\n"
				}
				message = message + fmt.Sprintf("%s: %s", discipline, markItem.Mark)
				continue
			}
		}
	}

	if len(message) == 0 {
		return "За сегодня нет оценок", nil
	}

	return message, nil
}

// getSummaryMarksRaw возвращает сырой ответ страницы getSummaryMarks
func getSummaryMarksRaw(login, password string) (io.ReadCloser, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, errors.Wrap(err, "не могу создать cookiejar")
	}
	client := http.Client{Jar: jar}

	loginResp, err := client.PostForm(formURL, url.Values{
		"login_login":    {login},
		"login_password": {password},
	})
	if err != nil {
		return nil, errors.Wrap(err, "не могу запостить форму")
	}
	defer loginResp.Body.Close()

	loginResponseBody, err := ioutil.ReadAll(loginResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Не могу прочитат тело ответа")
	}

	// может быть ошибка типа Account locked: too many login attempts. Please try again later
	if string(loginResponseBody) != `{"success": true,"redirect": "/"}` {
		return nil, errors.Errorf(
			"Неуспешный ответ страницы логина. Status:%s Body:%s",
			loginResp.Status,
			loginResponseBody,
		)
	}

	// после авторизации нужно зачем-то сходить на главную иначе не работает
	// ответ не смотрим, хз что там должно  быть и как это проверять
	_, err = client.Get(mainURL)
	if err != nil {
		return nil, errors.Wrap(err, "не могу сходить на главную")
	}

	summaryMarksRequest, err := http.NewRequest("GET", summaryMarksURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "не могу создать запрос на GetSummaryMarks")
	}
	q := summaryMarksRequest.URL.Query()
	q.Add("date", time.Now().Format("2006-01-02"))
	summaryMarksRequest.URL.RawQuery = q.Encode()
	marksResp, err := client.Do(summaryMarksRequest)
	if err != nil {
		return nil, errors.Wrap(err, "не могу сходить на GetSummaryMarks")
	}

	return marksResp.Body, nil
}

//func ()
