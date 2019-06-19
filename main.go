package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	configFlag = flag.String("config", "", "path to config file")
	logger     *zap.Logger
	config     *BotConfig
	err        error
	conn       *DBConn
	client     *twitch.Client
)

func init() {
	flag.Parse()
	logger, _ = zap.NewDevelopment()

	config, err = NewBotConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			config.M.Lock()
			err = config.Load(*configFlag)
			if err != nil {
				logger.Fatal("Can not reload config", zap.String("path", *configFlag), zap.String("error", err.Error()))
			}
			config.M.Unlock()

		}
	}()
}

func main() {
	go func() {
		conn, err = NewDBConn(config)
		if err != nil {
			logger.Fatal("Can not create database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
	}()

	accountName, err := config.GetAccountName()
	if err != nil {
		logger.Fatal("Can not get 'AccountName'", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	token, err := config.GetToken()
	if err != nil {
		logger.Fatal("Can not get 'AccountToken'", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	client = twitch.NewClient(accountName, token)
	client.TLS = false

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		_, err = conn.AddData(&msg)
		if err != nil {
			logger.Warn("Can not add data to database", zap.String("error", err.Error()))
		}
	})

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
