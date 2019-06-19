package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/gempir/go-twitch-irc"
	_ "github.com/lib/pq"
)

type DBConn struct {
	c *sql.DB
	m *sync.RWMutex
}

func NewDBConn(cfg *BotConfig) (*DBConn, error) {
	passwd, err := cfg.GetDBPassword()
	if err != nil {
		return nil, err
	}
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Address, cfg.Port, cfg.Username, passwd, cfg.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DBConn{
		c: db,
		m: new(sync.RWMutex),
	}, nil
}

func (c *DBConn) Close() error {
	return c.c.Close()
}

func (c *DBConn) AddData(msg *twitch.PrivateMessage) (sql.Result, error) {
	return c.c.Exec("CALL add_data($1, $2, $3, $4, $5, $6, $7, $8)",
		msg.Message, msg.ID, msg.Time, msg.Channel, msg.RoomID, msg.User.DisplayName, msg.User.ID,
		msg.Tags["turbo"]+msg.Tags["mod"]+msg.Tags["subscriber"])
}
