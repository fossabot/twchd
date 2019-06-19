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
	mu *sync.RWMutex
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
		mu: new(sync.RWMutex),
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

	b.mu.Lock()
	defer b.mu.Unlock()
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
	b.mu.RLock()
	defer b.mu.RUnlock()
	return getFieldOrEnv(b.AccountName)
}

// GetToken return token from environment or config file
func (b *BotConfig) GetToken() (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return getFieldOrEnv(b.AccountToken)
}
func (b *BotConfig) GetAccountsList() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.AccountsList
}

func (b *BotConfig) GetAddress() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Address
}

func (b *BotConfig) GetPort() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Port
}

func (b *BotConfig) GetUsername() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Username
}

func (b *BotConfig) GetDBPassword() (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return getFieldOrEnv(b.Password)
}

func (b *BotConfig) GetDatabase() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Database
}

func (b *BotConfig) Dump() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return fmt.Sprintf("%+v\n", b)
}
