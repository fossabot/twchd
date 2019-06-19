package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/lib/pq"
)

type DBConn struct {
	*sql.DB
	mu *sync.RWMutex
}

func NewDBConn(cfg *BotConfig) (*DBConn, error) {
	db, err := Connect(cfg)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		DB: db,
		mu: new(sync.RWMutex),
	}, nil
}

func Connect(cfg *BotConfig) (*sql.DB, error) {
	passwd, err := cfg.GetDBPassword()
	if err != nil {
		return nil, err
	}
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.GetAddress(), cfg.GetPort(), cfg.GetUsername(), passwd, cfg.GetDatabase())
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, err
}

func (c *DBConn) Reconnect(cfg *BotConfig) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	db, err := Connect(cfg)
	if err != nil {
		return err
	}
	c.Close()
	c.DB = db
	return nil
}

func (c *DBConn) AddData(msg *twitch.PrivateMessage) (sql.Result, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Exec("CALL add_data($1, $2, $3, $4, $5, $6, $7, $8)",
		msg.Message, msg.ID, msg.Time, msg.Channel, msg.RoomID, msg.User.DisplayName, msg.User.ID,
		msg.Tags["turbo"]+msg.Tags["mod"]+msg.Tags["subscriber"])
}
