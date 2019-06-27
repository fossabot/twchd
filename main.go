package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	configFlag = flag.String("config", "", "path to config file")
	logger     *zap.Logger
	config     *BotConfig
	err        error
	conn       *sql.DB
	initDBDone chan bool
	client     *twitch.Client
)

func main() {
	flag.Parse()
	logger, _ = zap.NewDevelopment()

	config, err = NewBotConfig(*configFlag)
	if err != nil {
		logger.Fatal("Can not load config", zap.String("path", *configFlag), zap.String("error", err.Error()))
	}
	logger.Info("Bot config loaded", zap.String("config", config.Dump()))

	initDBDone = make(chan bool)
	go func() {
		conn, err = NewDBConn(config)
		if err != nil {
			logger.Fatal("Can not create database connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
		initDBDone <- true
		logger.Info("DB connection created")
	}()

	client, err = NewTwitchClient(config)
	if err != nil {
		logger.Fatal("Can not create twitch client", zap.String("config", config.Dump()), zap.String("error", err.Error()))
	}
	logger.Info("Twitch client created")

	var queryStr = "CALL add_data($1, $2, $3, $4, $5, $6, $7, $8)"
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		var role = msg.Tags["turbo"] + msg.Tags["mod"] + msg.Tags["subscriber"]

		_, err = conn.Exec(queryStr, msg.ID, msg.Time, msg.RoomID, msg.Channel,
			msg.User.ID, msg.User.DisplayName, role, msg.Message)
		if err != nil {
			logger.Warn("Can not add data to database", zap.String("error", err.Error()))
		}
	})

	<-initDBDone
	logger.Info("Start connecting to twitch")
	err = client.Connect()
	if err != nil {
		logger.Warn("Error during connection to twitch", zap.String("error", err.Error()))
	}
}

func NewDBConn(cfg *BotConfig) (*sql.DB, error) {
	passwd, err := cfg.GetDBPassword()
	if err != nil {
		return nil, err
	}
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Address, cfg.Port, cfg.Username, passwd, cfg.Database)
	return sql.Open("postgres", connStr)
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
	for _, channel := range cfg.AccountsList {
		client.Join(channel)
	}
	return client, nil
}
