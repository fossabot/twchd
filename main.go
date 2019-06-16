package main

import (
	"flag"
	"log"

	"github.com/gempir/go-twitch-irc"
	"go.uber.org/zap"
)

var configFlag = flag.String("config", "", "path to config file")

func main() {
	flag.Parse()
	logger, _ := zap.NewDevelopment()

	config, err := NewBotConfig(*configFlag)
	if err != nil {
		log.Fatalln(err)
	}

	accountName, err := config.GetAccountName()
	if err != nil {
		log.Fatalln(err)
	}
	token, err := config.GetToken()
	if err != nil {
		log.Fatalln(err)
	}
	var client = twitch.NewClient(accountName, token)
	client.TLS = false

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {

	})

	for _, channel := range config.AccountsList {
		client.Join(channel)
	}
	if client.Connect() != nil {
		log.Fatalln(err)
	}
}
