package main

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
)

// VerifyPath verifies existance file
func VerifyPath(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}
	if !strings.HasSuffix(filename, ".yml") && !strings.HasSuffix(filename, ".yaml") {
		return errors.New("unsupported file format")
	}
	return nil
}

// BotConfig struct represent config from file
type BotConfig struct {
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
	Index        string   `yaml:"index"`
	Type         string   `yaml:"type"`
	Pipeline     string   `yaml:"pipeline"`
	Address      string   `yaml:"address"`
	Period       int      `yaml:"period"`
	Messages     int      `yaml:"n_messages"`
}

// NewBotConfig takes config file and return BotConfig struct
func NewBotConfig(filename string) (config *BotConfig, err error) {
	err = VerifyPath(filename)
	if err != nil {
		return nil, err
	}
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config = new(BotConfig)
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return nil, err
	}
	return
}

func getFieldOrEnv(field string) (string, error) {
	if field == "" {
		return "", errors.New("Config field empty")
	}
	pattern, err := regexp.Compile(`^\${(.*?)}$`)
	if err != nil {
		return "", err
	}
	var textVar = pattern.FindStringSubmatch(field)
	if textVar == nil {
		return field, nil
	}
	var envVar = os.Getenv(textVar[1])
	if envVar == "" {
		return "", errors.New("Environment variable does not set")
	}
	return envVar, nil
}

// GetAccountName return account name from environment or config file
func (b *BotConfig) GetAccountName() (string, error) {
	accountName, err := getFieldOrEnv(b.AccountName)
	if err != nil {
		return "", err
	}
	return accountName, nil
}

// GetToken return token from environment or config file
func (b *BotConfig) GetToken() (string, error) {
	token, err := getFieldOrEnv(b.AccountToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (b *BotConfig) GetPeriod() time.Duration {
	return time.Duration(b.Period)
}
