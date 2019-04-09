package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gempir/go-twitch-irc"
)

// Check unrecoverable error and panic
func Check(err error) {
	if err != nil {
		panic(err)
	}
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

// Dump whole message for debug purpose
func (m *Message) String() string {
	str, err := json.Marshal(m)
	Check(err)

	return string(str)
}

// selectChannelName extract channel name from raw string in message
func selectChannelName(raw string) (string, error) {
	re, err := regexp.Compile(`PRIVMSG\s#(.*?)\s`)
	Check(err)

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
	Check(flags.VerifyPath())

	config := NewBotConfig(flags.ConfigPath)

	logger := log.New(os.Stderr, "", 0)

	twClient := twitch.NewClient(config.AccountName, config.AccountToken)
	twClient.TLS = false

	twClient.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		esMsg, err := NewMessage(&message)
		if err != nil {
			return
		}

		if flags.DebugOutput {
			logger.Println(esMsg)
		}
	})

	config.JoinAllTo(twClient)

	Check(twClient.Connect())

}
