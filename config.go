package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	TelegramBot struct {
		Token   string   `yaml:"token"`
		ChatIds []string `yaml:"chat_ids"`
		Proxy   string   `yaml:"proxy"`
	} `yaml:"telegram_bot"`

	Credentials struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
	}
}

func getConfig(fileName string) (config, error) {
	conf := config{}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return conf, fmt.Errorf("can't read file %s, %s", fileName, err)
	}

	if err := yaml.Unmarshal(file, &conf); err != nil {
		return conf, fmt.Errorf("can't unmarshal config file: %v", err)
	}

	return conf, nil
}
