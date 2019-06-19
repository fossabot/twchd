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

func main() {
	flag.Parse()
	logger, _ = zap.NewDevelopment()

	config, err = NewBotConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}

	go func() {
		conn, err = NewDBConn(config)
		if err != nil {
			logger.Fatal("Can not create database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
	}()

	client, err = NewTwitchClient(config)
	if err != nil {
		logger.Fatal("Can not create twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			err = config.Load(*configFlag)
			if err != nil {
				logger.Fatal("Can not reload config", zap.String("path", *configFlag), zap.String("error", err.Error()))
			}

			err = conn.Reconnect(config)
			if err != nil {
				logger.Fatal("Can not create new connection to database", zap.String("config", config.Dump()), zap.String("error", err.Error()))
			}

			client, err = NewTwitchClient(config)
			if err != nil {
				logger.Fatal("Can not recreate twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
			}
		}
	}()

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		_, err = conn.AddData(&msg)
		if err != nil {
			logger.Warn("Can not add data to database", zap.String("error", err.Error()))
		}
	})

	RetryConnect(client, logger, 10, 6)
}

func NewTwitchClient(cfg *BotConfig) (*twitch.Client, error) {
	accountName, err := cfg.GetAccountName()
	if err != nil {
		return nil, err
	}
	token, err := cfg.GetToken()
	if err != nil {
		return nil, err
	}
	var client = twitch.NewClient(accountName, token)
	client.TLS = false

	for _, channel := range cfg.GetAccountsList() {
		client.Join(channel)
	}

	return client, nil
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
