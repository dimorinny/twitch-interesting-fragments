package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-chat-api"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	irc "github.com/fluffle/goirc/client"
	"log"
)

var (
	config     configuration.Configuration
	connection *irc.Conn
)

func initConfiguration() {
	config = configuration.Configuration{}

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
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
	initTwitchConnection()
}

func main() {
	twitchConfiguration := twitchchat.NewConfiguration(
		config.Nickname,
		config.Oauth,
		config.Channel,
	)

	chat := twitchchat.NewChat(
		twitchConfiguration,
		connection,
	)

	runWithChannels(chat)
}

func runWithChannels(twitch *twitchchat.Chat) {
	disconnected := make(chan struct{})
	connected := make(chan struct{})
	errStream := make(chan error)
	message := make(chan string)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
			case <-connected:
				fmt.Println("Connected")
			case err := <-errStream:
				fmt.Println(err)
			case newMessage := <-message:
				fmt.Println(newMessage)
			}
		}
	}()

	twitch.ConnectWithChannels(connected, disconnected, errStream, message)
}

func runWithCallbacks(twitch *twitchchat.Chat) {
	twitch.ConnectWithCallbacks(
		func() {
			fmt.Println("Connected")
		},
		func() {
			fmt.Println("Disconnected")
		},
		func(err error) {
			fmt.Println(err)
		},
		func(message string) {
			fmt.Println(message)
		},
	)
}
