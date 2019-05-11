package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/olivere/elastic"
)

func main() {
	flags := NewFlagsCLI()
	config, err := NewBotConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	esClient, err := elastic.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	var esCtx = context.Background()

	defer esClient.CloseIndex(config.Index).Do(esCtx)
	defer esClient.Flush(config.Index).Do(esCtx)

	ingestPipeline, err := Asset("assets/pipeline.json")
	if err != nil {
		log.Fatalln(err)
	}
	var newPipeline = strings.Replace(string(ingestPipeline), "${TWITCH_CHANNEL}", config.Index, 1)

	_, err = esClient.IngestPutPipeline(config.Pipeline).
		BodyString(newPipeline).
		Do(esCtx)
	if err != nil {
		log.Fatalln(err)
	}
	accountName, err := config.GetAccountName()
	if err != nil {
		log.Fatalln(err)
	}
	token, err := config.GetToken()
	if err != nil {
		log.Fatalln(err)
	}
	var twClient = twitch.NewClient(accountName, token)
	twClient.TLS = false

	twClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		var builder strings.Builder
		fmt.Fprintf(&builder, "{\"msg\": \"%s,%s,", message.RoomID, message.Channel)
		fmt.Fprintf(&builder, "%s,%d,", message.Message, message.Time.Unix())
		fmt.Fprintf(&builder, "%s,%s,", message.User.DisplayName, message.User.ID)
		fmt.Fprintf(&builder, "%s,%s,%s\"}", message.Tags["turbo"], message.Tags["subscriber"], message.Tags["mod"])
		var message2 = builder.String()

		resp, err := esClient.Index().
			Index(config.Index).
			Type(config.Type).
			Pipeline(config.Pipeline).
			BodyString(message2).
			Do(esCtx)
		if err != nil {
			return
		}

		if flags.DebugOutput {
			fmt.Fprintf(os.Stderr, "[%v] %v: %v\n", resp.Id, message.User.DisplayName, message.Message)
		}
	})

	for _, channel := range config.AccountsList {
		twClient.Join(channel)
	}
	if twClient.Connect() != nil {
		log.Fatalln(err)
	}
}
