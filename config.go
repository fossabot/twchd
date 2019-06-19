package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/go-yaml/yaml"
)

// VerifyConfigPath verifies existance file
func VerifyConfigPath(filename string) error {
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
	M *sync.RWMutex
	// Twitch Client settings
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
	// Postgres DB settings
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// NewBotConfig takes config file and return BotConfig struct
func NewBotConfig(filename string) (*BotConfig, error) {
	var config = &BotConfig{
		M: new(sync.Mutex),
	}
	var err = config.Load(filename)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (b *BotConfig) Load(filename string) error {
	var err = VerifyConfigPath(filename)
	if err != nil {
		return err
	}
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(rawConfig, b)
	if err != nil {
		return err
	}
	return nil
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
	return getFieldOrEnv(b.AccountName)
}

// GetToken return token from environment or config file
func (b *BotConfig) GetToken() (string, error) {
	return getFieldOrEnv(b.AccountToken)
}

func (b *BotConfig) GetDBPassword() (string, error) {
	return getFieldOrEnv(b.Password)
}

func (b *BotConfig) Dump() string {
	return fmt.Sprintf("%+v\n", b)
}
