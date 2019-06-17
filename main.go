package main

import (
	"database/sql"
	"flag"
	"fmt"

	"time"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/lib/pq"
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

	var db *sql.DB
	var initDBDone = make(chan bool)
	defer db.Close()
	go func() {
		DBPassword, err := config.GetDBPassword()
		if err != nil {
			logger.Fatal("Can not get 'GetDBPassword'", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
		var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Address, config.Port, config.Username, DBPassword, config.Database)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			logger.Fatal("Can not open postgres connection", zap.String("config", config.Dump()), zap.String("error", err.Error()))
		}
		initDBDone <- true
	}()

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

	<-initDBDone
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		var queryStr = "CALL add_data($1, $2, $3, $4, $5, $6, $7, $8::BIT(3))"
		var role = msg.Tags["turbo"] + msg.Tags["mod"] + msg.Tags["subscriber"]

		_, err := db.Query(queryStr, msg.Message, msg.ID, msg.Time, msg.Channel, msg.RoomID, msg.User.DisplayName, msg.User.ID, role)
		if err != nil {
			logger.Warn("Can not make query to database", zap.String("error", err.Error()))
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
