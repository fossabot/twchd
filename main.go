package main

import (
	"flag"

	"github.com/gempir/go-twitch-irc"
	"go.uber.org/zap"
)

var (
	configFlag = flag.String("config", "", "path to config file")
)

func main() {
	flag.Parse()
	logger, _ := zap.NewDevelopment()

	config, err := NewBotConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}

	accountName, err := config.GetAccountName()
	if err != nil {
		logger.Fatal("Can not get 'AccountName'", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	token, err := config.GetToken()
	if err != nil {
		logger.Fatal("Can not get 'AccountToken'", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	var client = twitch.NewClient(accountName, token)
	client.TLS = false

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {

	})

	for _, channel := range config.AccountsList {
		client.Join(channel)
	}
	if client.Connect() != nil {
		logger.Fatal("Error during twitch connection")
	}
}
