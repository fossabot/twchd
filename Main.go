package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gempir/go-twitch-irc"
	"github.com/olivere/elastic"
)

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

// NewMessage select fields from twitch.Message
func NewMessage(msg *twitch.Message) (*Message, error) {
	sentTS, err := strconv.Atoi(msg.Tags["tmi-sent-ts"])
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(msg.Tags["user-id"])
	if err != nil {
		return nil, err
	}

	mod, err := strIntToBool(msg.Tags["mod"])
	if err != nil {
		return nil, err
	}

	sub, err := strIntToBool(msg.Tags["subscriber"])
	if err != nil {
		return nil, err
	}

	turbo, err := strIntToBool(msg.Tags["turbo"])
	if err != nil {
		return nil, err
	}

	channelID, err := strconv.Atoi(msg.ChannelID)
	if err != nil {
		return nil, err
	}

	channelName, err := selectChannelName(msg.Raw)
	if err != nil {
		return nil, err
	}

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
			ChannelID:   uint64(channelID),
			ChannelName: channelName,
		},
	}, nil
}

// Dump whole message
func (m *Message) Dump() (d []byte) {
	d, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return
}

// selectChannelName extract channel name from raw string in message
func selectChannelName(raw string) (string, error) {
	re, err := regexp.Compile(`PRIVMSG\s#(.*?)\s`)
	if err != nil {
		panic(err)
	}

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

func main() {
	flags := NewFlagsCLI()
	err := flags.VerifyPath()
	if err != nil {
		log.Fatalln(err)
	}
	config := NewBotConfig(flags.ConfigPath)

	logger := log.New(os.Stderr, "", 0)

	esClient, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	exists, err := esClient.IndexExists(config.IndexES).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exists {
		logger.Printf("index '%v' does not exists, creating...\n", config.IndexES)

		mapping, err := ioutil.ReadFile("mapping.json")
		if err != nil {
			panic(err)
		}

		createIndex, err := esClient.CreateIndex(config.IndexES).Body(string(mapping)).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic(errors.New("index '" + config.IndexES + "' does not created"))
		}
	}

	twClient := twitch.NewClient(config.AccountName, config.AccountToken)
	twClient.TLS = false

	twClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		esMsg, err := NewMessage(&message)
		if err != nil {
			return
		}

		_, err = esClient.Index().
			Index(config.IndexES).
			Type(config.TypeES).
			BodyJson(esMsg).
			Do(context.Background())
		if err != nil {
			return
		}

		if flags.DebugOutput {
			logger.Println(*esMsg)
		}
	})

	config.JoinAllTo(twClient)

	if twClient.Connect() != nil {
		panic(err)
	}

	_, err = esClient.Flush().Index(config.IndexES).Do(context.Background())
	if err != nil {
		panic(err)
	}
}
