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
	client     *TwitchClient
)

func main() {
	flag.Parse()
	logger, _ = zap.NewDevelopment()

	config, err = NewBotConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}
	logger.Info("Bot config loaded", zap.String("config", config.Dump()))

	go func() {
		conn, err = NewDBConn(config)
		if err != nil {
			logger.Fatal("Can not create database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
		logger.Info("DB connection created")
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)

	go func() {
		logger.Info("Starting handling SIGHUP")
		for {
			<-s
			err = config.Load(*configFlag)
			if err != nil {
				logger.Fatal("Can not reload config", zap.String("path", *configFlag), zap.String("error", err.Error()))
			}
			logger.Info("Bot config reloaded")

			go func() {
				err = conn.Reconnect(config)
				if err != nil {
					logger.Fatal("Can not update database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
				}
				logger.Info("DB connection updated")
			}()

			client.Disconnect()
			err = client.Reconfigure(config)
			client.OnPrivateMessage(addDataHandler)
			if err != nil {
				logger.Fatal("Can not recreate twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
			}
			logger.Info("Twitch client updated")
		}
	}()

	client, err = NewTwitchClient(config)
	if err != nil {
		logger.Fatal("Can not create twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	logger.Info("Twitch client created")

	client.OnPrivateMessage(addDataHandler)

	logger.Info("Start connecting to twitch")
	var ticker = time.NewTicker(time.Duration(10) * time.Second)
	for {
		err = client.Connect()
		if err != nil {
			logger.Warn("Error during connection to twitch", zap.String("error", err.Error()))
		}
		<-ticker.C
	}
}

func addDataHandler(msg twitch.PrivateMessage) {
	_, err = conn.AddData(&msg)
	if err != nil {
		logger.Warn("Can not add data to database", zap.String("error", err.Error()))
	}
}
