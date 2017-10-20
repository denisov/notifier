package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const formURL = "http://billing.kengudetyam.ru/cabinet/Account/Login"

// getContent возвращает контент адмики
func getContent(login, password string) (*http.Response, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, fmt.Errorf("не могу создать cookiejar %v", err)
	}
	client := http.Client{Jar: jar}
	resp, err := client.PostForm(formURL, url.Values{
		"UserName": {login},
		"Password": {password},
	})
	if err != nil {
		return nil, fmt.Errorf("не могу запостить форму %v", err)
	}

	return resp, nil
}

// getBalance возвращает значение баданса из html тэга
func getBalance(pageResponse *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(pageResponse)
	if err != nil {
		return "", fmt.Errorf("не могу создать документ для парсинга URL: %v", err)
	}

	balanceSelection := doc.Find(".balance")
	if balanceSelection.Length() == 0 {
		return "", fmt.Errorf("не могу найти '.balance' в html коде")
	}

	return balanceSelection.First().Text(), nil
}
