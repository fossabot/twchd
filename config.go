package main

import (
	"io/ioutil"

	"github.com/gempir/go-twitch-irc"
	"github.com/go-yaml/yaml"
)

// BotConfig struct represent config from file
type BotConfig struct {
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
	LoginMB      string   `yaml:"login"`
	PasswdMB     string   `yaml:"password"`
	AddrMB       string   `yaml:"address"`
	PortMB       string   `yaml:"port"`
	Exchange     string   `yaml:"exchange_name"`
}

// NewBotConfig takes config file and return BotConfig struct
func NewBotConfig(filename string) (config *BotConfig) {
	rawConfig, err := ioutil.ReadFile(filename)
	Check(err)

	config = new(BotConfig)
	Check(yaml.Unmarshal(rawConfig, config))
	return
}

// JoinAllTo joins client to all accounts from config
func (c *BotConfig) JoinAllTo(client *twitch.Client) {
	for _, channel := range c.AccountsList {
		client.Join(channel)
	}
}
