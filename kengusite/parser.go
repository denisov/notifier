package kengusite

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

// Parser это парсер сайта кенгудетям
type Parser struct {
	login    string
	password string
}

// NewParser создаёт парсер сайта кенгудетям
func NewParser(login, password string) Parser {
	return Parser{
		login:    login,
		password: password,
	}
}

const formURL = "http://billing.kengudetyam.ru/cabinet/Account/Login"

// getContent возвращает контент адмики
func (p Parser) getContent() (*http.Response, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, errors.Wrap(err, "не могу создать cookiejar")
	}
	client := http.Client{Jar: jar}
	resp, err := client.PostForm(formURL, url.Values{
		"UserName": {p.login},
		"Password": {p.password},
	})
	if err != nil {
		return nil, errors.Wrap(err, "не могу запостить форму")
	}

	return resp, nil
}

// GetData возвращает баланс
func (p Parser) GetData() (string, error) {
	pageResponse, err := p.getContent()
	if err != nil {
		return "", errors.WithStack(err)
	}

	doc, err := goquery.NewDocumentFromResponse(pageResponse)
	if err != nil {
		return "", errors.Wrap(err, "не могу создать документ для парсинга URL")
	}

	balanceSelection := doc.Find(".balance")
	if balanceSelection.Length() == 0 {
		return "", errors.New("не могу найти '.balance' в html коде")
	}

	balance := balanceSelection.First().Text()
	// заменяем неразрывный пробел на нормальный
	balance = strings.Replace(balance, "\u00a0", " ", -1)

	return strings.Replace(balanceSelection.First().Text(), "\u00a0", " ", -1), nil
}
