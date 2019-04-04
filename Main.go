package main

import (
	"errors"
	"flag"
	"fmt"
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

// JoinAllTo joins client to all accounts from config
func (c *BotConfig) JoinAllTo(client *twitch.Client) {
	for _, channel := range c.AccountsList {
		client.Join(channel)
	}
}

// BotCliFlags store cli flags after parse
type BotCliFlags struct {
	ConfigPath  string
	DebugOutput bool
}

// NewCliFlags parse cli args and return BotCliFlags struct
func NewCliFlags() *BotCliFlags {
	flagConfig := flag.String("config", "", "path to config file")
	flagDebug := flag.Bool("debug", false, "addition output to syslog")
	flag.Parse()

	return &BotCliFlags{
		ConfigPath:  *flagConfig,
		DebugOutput: *flagDebug,
	}
}

// VerifyPath verifies ConfigPath
func (f *BotCliFlags) VerifyPath() error {
	if len(f.ConfigPath) == 0 {
		return errors.New("path to config file does not passed")
	}

	if _, err := os.Stat(f.ConfigPath); os.IsNotExist(err) {
		return errors.New("file does not exists")
	}

	return nil
}

// TODO: EmoteFilter

func main() {
	flags := NewCliFlags()
	check(flags.VerifyPath())

	// journal, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_DAEMON, "botbot.com")
	// check(err)

	config := ConfigParser(flags.ConfigPath)

	twClient := twitch.NewClient(config.AccountName, config.AccountToken)

	// esClient, err := elastic.NewClient()
	// check(err)

	twClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		fmt.Printf("%+v\n", message)
		if flags.DebugOutput {
			//journal.Debug()
		}
	})

	config.JoinAllTo(twClient)

	check(twClient.Connect())
}
