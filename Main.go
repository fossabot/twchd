package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"regexp"
	"strconv"

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

// User represent part related with user info
type User struct {
	DisplayName string `json:"display_name"`
	UserID      uint64 `json:"user_id"`
	IsMod       bool   `json:"moderator"`
	IsSub       bool   `json:"subscriber"`
	IsTurbo     bool   `json:"turbo"`
}

// Channel represent part related with channel info
type Channel struct {
	ChannelName string `json:"channel_name"`
	ChannelID   uint64 `json:"channel_id"`
}

// Message maps go-struct to elasticsearch document
type Message struct {
	Text      string `json:"text"`
	TimeEpoch uint64 `json:"time"`
	User      `json:"user"`
	Channel   `json:"channel"`
}

// MakeMessage select fields from twitch.Message
func MakeMessage(msg *twitch.Message) *Message {
	sentTS, err := strconv.Atoi(msg.Tags["tmi-sent-ts"])
	check(err)

	userID, err := strconv.Atoi(msg.Tags["user-id"])
	check(err)

	mod, err := strIntToBool(msg.Tags["mod"])
	check(err)

	sub, err := strIntToBool(msg.Tags["subscriber"])
	check(err)

	turbo, err := strIntToBool(msg.Tags["turbo"])
	check(err)

	chID, err := strconv.Atoi(msg.ChannelID)
	check(err)

	chName, err := selectChannelName(msg.Raw)
	check(err)

	return &Message{
		Text:      msg.Text,
		TimeEpoch: uint64(sentTS),
		User: User{
			DisplayName: msg.Tags["display-name"],
			UserID:      uint64(userID),
			IsMod:       mod,
			IsSub:       sub,
			IsTurbo:     turbo,
		},
		Channel: Channel{
			ChannelID:   uint64(chID),
			ChannelName: chName,
		},
	}
}

func (m *Message) Dump() string {
	str, err := json.Marshal(&m)
	check(err)

	return string(str)
}

// selectChannelName extract channel name from raw string in message
func selectChannelName(raw string) (string, error) {
	re, err := regexp.Compile(`PRIVMSG\s#(.*?)\s`)
	check(err)

	match := re.FindStringSubmatch(raw)
	if len(match) == 0 {
		return "", errors.New("channel not found")
	}
	return match[1], nil
}

// strIntToBool convert integer as string to bool
func strIntToBool(str string) (bool, error) {
	if i, err := strconv.Atoi(str); err == nil {
		return i != 0, nil
	}
	return false, errors.New("Can not convert")
}

// TODO: DO NOT EXIT ON MESSAGE PROCESSING ERROR
func main() {
	flags := NewCliFlags()
	check(flags.VerifyPath())

	journal, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_DAEMON, "botbot.com")
	check(err)

	config := ConfigParser(flags.ConfigPath)

	twClient := twitch.NewClient(config.AccountName, config.AccountToken)

	// ctx := context.Background()
	// esClient, err := elastic.NewClient()
	// check(err)

	twClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		esMsg := MakeMessage(&message)
		if flags.DebugOutput {
			journal.Debug(esMsg.Dump())
		}
	})

	config.JoinAllTo(twClient)

	check(twClient.Connect())
}
