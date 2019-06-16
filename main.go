package main

import (
	"flag"
	"time"

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

	// client.OnPrivateMessage()

	for _, channel := range config.AccountsList {
		client.Join(channel)
	}
	RetryConnect(client, logger, 10, 6)
}

func RetryConnect(client *twitch.Client, logger *zap.Logger, period int, attempts int) {
	var ticker = time.NewTicker(time.Duration(period) * time.Second)
	for attempt := 0; attempt < attempts; attempt++ {
		if client.Connect() != nil {
			logger.Warn("Error during connection to twitch", zap.Int("attempt", attempt), zap.Int("timeout", period))
		}
		<-ticker.C
	}
	ticker.Stop()
	logger.Fatal("Error during connection to twitch. Attempts exceeded")
}
