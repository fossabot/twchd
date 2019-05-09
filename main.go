package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/olivere/elastic"
)

func extractFile(filename string) []byte {
	file, err := assets.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	text, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	return text
}

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
	var esCtx = context.Background()

	defer esClient.CloseIndex(config.Index).Do(esCtx)
	defer esClient.Flush(config.Index).Do(esCtx)

	_, err = esClient.IngestPutPipeline(config.Pipeline).
		BodyString(string(extractFile("/pipeline.json"))).
		Do(esCtx)
	if err != nil {
		log.Fatalln(err)
	}

	var twClient = twitch.NewClient(config.GetAccountName(), config.GetToken())
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
			fmt.Fprintln(os.Stderr, resp.Id, resp.Result, message.User.DisplayName, message.Message)
		}
	})

	for _, channel := range config.AccountsList {
		twClient.Join(channel)
	}
	if twClient.Connect() != nil {
		log.Fatalln(err)
	}
}
