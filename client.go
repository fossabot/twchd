package main

import (
	"sync"

	"github.com/gempir/go-twitch-irc"
)

type TwitchClient struct {
	*twitch.Client
	mu *sync.RWMutex
}

func NewTwitchClient(cfg *BotConfig) (*TwitchClient, error) {
	var client = &TwitchClient{
		Client: nil,
		mu:     new(sync.RWMutex),
	}
	var err = client.Reconfigure(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *TwitchClient) Reconfigure(cfg *BotConfig) error {
	accountName, err := cfg.GetAccountName()
	if err != nil {
		return err
	}
	token, err := cfg.GetToken()
	if err != nil {
		return err
	}

	var client = twitch.NewClient(accountName, token)
	client.TLS = false
	for _, channel := range cfg.GetAccountsList() {
		client.Join(channel)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Client = client
	return nil
}

func (c *TwitchClient) Connect() error {
	c.mu.Lock()
	var err = c.Client.Connect()
	c.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}
