package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/olivere/elastic"
)

func main() {
	flags := NewFlagsCLI()
	config, err := NewBotConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	esClient, err := elastic.NewClient(elastic.SetURL(config.Address), elastic.SetSniff(false))
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

	var bulker = esClient.Bulk()

	twClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		var req = elastic.NewBulkIndexRequest().
			Index(config.Index).
			Type(config.Type).
			Pipeline(config.Pipeline).
			Doc(MakeRawMsg(&message))
		bulker.Add(req)
	})

	go func() {
		var ticker = time.NewTicker(config.GetPeriod() * time.Second)
		for {
			select {
			case <-ticker.C:
				bulker.Do(esCtx)
			}
		}
	}()

	for _, channel := range config.AccountsList {
		twClient.Join(channel)
	}
	if twClient.Connect() != nil {
		log.Fatalln(err)
	}
}

func MakeRawMsg(message *twitch.PrivateMessage) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "{\"msg\": \"%s,%s,", message.RoomID, message.Channel)
	fmt.Fprintf(&builder, "%s,%d,", message.Message, message.Time.Unix())
	fmt.Fprintf(&builder, "%s,%s,", message.User.DisplayName, message.User.ID)
	fmt.Fprintf(&builder, "%s,%s,%s\"}", message.Tags["turbo"], message.Tags["subscriber"], message.Tags["mod"])
	return builder.String()
}
