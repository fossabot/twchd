package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/olivere/elastic"
)

func main() {
	flags := NewFlagsCLI()
	err := flags.VerifyPath()
	if err != nil {
		log.Fatalln(err)
	}
	config := NewBotConfig(flags.ConfigPath)

	esClient, err := elastic.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	twClient := twitch.NewClient(config.GetAccountName(), config.GetToken())
	twClient.TLS = false

	twClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		var builder strings.Builder
		fmt.Fprintf(&builder, "{\"msg\": \"%s,%s,", message.RoomID, message.Channel)
		fmt.Fprintf(&builder, "%s,%d,", message.Message, message.Time.Unix())
		fmt.Fprintf(&builder, "%s,%s,", message.User.DisplayName, message.User.ID)
		fmt.Fprintf(&builder, "%s,%s,%s\"}", message.Tags["turbo"], message.Tags["subscriber"], message.Tags["mod"])
		var message2 = builder.String()

		resp, err := esClient.Index().
			Index(config.IndexES).
			Type(config.TypeES).
			Pipeline(config.PipelineES).
			BodyString(message2).
			Do(context.Background())
		if err != nil {
			return
		}

		if flags.DebugOutput {
			log.Println(resp.Id, resp.Result, message2)
		}
	})

	for _, channel := range config.AccountsList {
		twClient.Join(channel)
	}

	if twClient.Connect() != nil {
		panic(err)
	}
}
