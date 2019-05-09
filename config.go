package main

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// BotConfig struct represent config from file
type BotConfig struct {
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
	IndexES      string   `yaml:"index"`
	TypeES       string   `yaml:"type"`
	PipelineES   string   `yaml:"pipeline"`
}

// NewBotConfig takes config file and return BotConfig struct
func NewBotConfig(filename string) (config *BotConfig) {
	rawConfig, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config = new(BotConfig)
	if yaml.Unmarshal(rawConfig, config) != nil {
		panic(err)
	}
	return
}
