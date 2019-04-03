package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/gempir/go-twitch-irc"
	"github.com/go-yaml/yaml"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// BotConfig struct represent config from file
type BotConfig struct {
	AccountName  string   `yaml:"account_name"`
	AccountToken string   `yaml:"account_token"`
	AccountsList []string `yaml:"join_to"`
}

// ConfigParser takes config file and return BotConfig struct
func ConfigParser(filename string) (config *BotConfig) {
	rawConfig, err := ioutil.ReadFile(filename)
	check(err)

	config = new(BotConfig)
	check(yaml.Unmarshal(rawConfig, config))
	return
}

// NewClient creates twitch client from config
func (c *BotConfig) NewClient() *twitch.Client {
	return twitch.NewClient(c.AccountName, c.AccountToken)
}

// JoinAll joins client to all accounts from config
func (c *BotConfig) JoinAll(client *twitch.Client) {
	for _, channel := range c.AccountsList {
		client.Join(channel)
	}
}

var flagConfigPath string

// verify config path
func configPath() error {
	if len(flagConfigPath) == 0 {
		return errors.New("path to config file does not passed")
	}

	if _, err := os.Stat(flagConfigPath); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}

	return nil
}

func init() {
	const usageFlagConfig = "path to config file"
	flag.StringVar(&flagConfigPath, "c", "", usageFlagConfig)
	flag.StringVar(&flagConfigPath, "config", "", usageFlagConfig)
	flag.Parse()
}

func main() {
	check(configPath())
	config := ConfigParser(flagConfigPath)

	client := config.NewClient()
	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		log.Printf("%+v\n", user)
	})

	config.JoinAll(client)

	check(client.Connect())
}
