package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-chat-api"
	"github.com/dimorinny/twitch-interesting-fragments/api"
	"github.com/dimorinny/twitch-interesting-fragments/buffer"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	irc "github.com/fluffle/goirc/client"
	"log"
	"net/http"
	"time"
)

var (
	config     configuration.Configuration
	connection *irc.Conn
	uploader   *api.Uploader
)

func initConfiguration() {
	config = configuration.Configuration{}

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func initUploader() {
	uploader = api.NewUploader(config, http.DefaultClient)
}

func initTwitchIrcConfig() (ircConfig *irc.Config) {
	ircConfig = irc.NewConfig(config.Nickname)
	ircConfig.Server = config.Host
	ircConfig.Pass = config.Oauth
	return
}

func initTwitchConnection() {
	connection = irc.Client(initTwitchIrcConfig())
	return
}

func init() {
	initConfiguration()
	initUploader()
	initTwitchConnection()
}

func main() {
	uploadExample()
}

func uploadExample() {
	result, err := uploader.Upload()

	if err != nil {
		log.Fatal(err)
	}

	print(result.Data.Url)
}

func bufferExample() {
	twitchConfiguration := twitchchat.NewConfiguration(
		config.Nickname,
		config.Oauth,
		config.Channel,
	)

	chat := twitchchat.NewChat(
		twitchConfiguration,
	)

	message := make(chan string)

	timeBuffer := buffer.NewMessagesBuffer(message, time.Second*5)

	go ircChatExample(chat, message)

	bufferedChannel := timeBuffer.Start()

	go func() {
		time.Sleep(15 * time.Second)
		fmt.Println("Stopped1")
		timeBuffer.Stop()
	}()

	for {
		items, ok := <-bufferedChannel

		if !ok {
			break
		}

		fmt.Println(items)
		fmt.Println(ok)
	}

	bufferedChannel = timeBuffer.Start()

	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("Stopped2")
		timeBuffer.Stop()
	}()

	for {
		items, ok := <-bufferedChannel

		if !ok {
			break
		}

		fmt.Println(items)
		fmt.Println(ok)
	}
}

func ircChatExample(twitch *twitchchat.Chat, message chan string) {
	disconnected := make(chan struct{})
	connected := make(chan struct{})
	errStream := make(chan error)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
			case <-connected:
				fmt.Println("Connected")
			case err := <-errStream:
				fmt.Println(err)
			case _ = <-message:
			}
		}
	}()

	twitch.ConnectWithChannels(connected, disconnected, errStream, message)
}
