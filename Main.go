package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"
)

func main() {
	client := twitch.NewClient("vanya109", "oauth:kic4lme1k221w4c6soyrppnwros6b0")
	client.OnNewMessage(func(channel string, user twitch.User, message twitch.Message) {
		fmt.Println(user.DisplayName, message.Text)
	})

	client.Join("vanya83")

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
