package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

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

	logger := log.New(os.Stderr, "", 0)

	esClient, err := elastic.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	inxCtx := context.Background()
	exists, err := esClient.IndexExists(config.IndexES).Do(inxCtx)
	if err != nil {
		log.Fatalln(err)
	}
	if !exists {
		logger.Printf("index '%v' does not exists, creating...\n", config.IndexES)

		file, err := assets.Open("/mapping.json")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		mapping, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalln(err)
		}

		createIndex, err := esClient.CreateIndex(config.IndexES).Body(string(mapping)).Do(inxCtx)
		if err != nil {
			log.Fatalln(err)
		}
		if !createIndex.Acknowledged {
			log.Fatalln("index '" + config.IndexES + "' does not created")
		}
	}

	twClient := twitch.NewClient(config.AccountName, config.AccountToken)
	twClient.TLS = false

	twClient.OnPrivateMessage(func(message twitch.PrivateMessage) {

		// message.RoomID
		// message.Channel

		// message.Message
		// message.Time

		// message.User.DisplayName
		// message.User.ID
		// message.Tags["turbo"]
		// message.Tags["subscriber"]
		// message.Tags["mod"]

		// _, err = esClient.Index().
		// 	Index(config.IndexES).
		// 	Type(config.TypeES).
		// 	BodyString().
		// 	Do(context.Background())
		// if err != nil {
		// 	return
		// }

		if flags.DebugOutput {
			logger.Printf("%+v\n", message)
		}
	})

	for _, channel := range config.AccountsList {
		twClient.Join(channel)
	}

	if twClient.Connect() != nil {
		panic(err)
	}

	_, err = esClient.Flush().Index(config.IndexES).Do(inxCtx)
	if err != nil {
		panic(err)
	}
}
