package main

import (
	"io/ioutil"
	"log"

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

// JoinAll joins client to all accounts from config
func (c *BotConfig) JoinAll(client *twitch.Client) {
	for _, channel := range c.AccountsList {
		client.Join(channel)
	}
}

// TODO: flags
// -c to pass config file
func main() {
	config := ConfigParser("config.yml")

	client := twitch.NewClient(config.AccountName, config.AccountToken)
	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		log.Printf("%+v\n", user)
	})

	config.JoinAll(client)

	check(client.Connect())
}
