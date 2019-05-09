package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/go-yaml/yaml"
)

// BotConfig struct represent config from file
type BotConfig struct {
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
	Index        string   `yaml:"index"`
	Type         string   `yaml:"type"`
	Pipeline     string   `yaml:"pipeline"`
}

// NewBotConfig takes config file and return BotConfig struct
func NewBotConfig(filename string) (config *BotConfig) {
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	config = new(BotConfig)
	if yaml.Unmarshal(rawConfig, config) != nil {
		log.Fatalln(err)
	}
	return
}

func getFieldOrEnv(field string) (string, error) {
	pattern, err := regexp.Compile(`^\${(.*?)}$`)
	if err != nil {
		log.Fatalln(err)
	}
	var envVar = pattern.FindStringSubmatch(field)
	if envVar == nil {
		return field, nil
	}

	var env = os.Getenv(envVar[1])
	if env == "" {
		return "", errors.New("Environment variable does not set")
	}
	return env, nil
}

// GetAccountName return account name from environment or config file
func (b *BotConfig) GetAccountName() string {
	accountName, err := getFieldOrEnv(b.AccountName)
	if err != nil {
		log.Fatalln(err)
	}
	return accountName
}

// GetToken return token from environment or config file
func (b *BotConfig) GetToken() string {
	token, err := getFieldOrEnv(b.AccountToken)
	if err != nil {
		log.Fatalln(err)
	}
	return token
}
