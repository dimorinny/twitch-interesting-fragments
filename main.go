package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-interesting-fragments/chat"
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
	client := chat.NewClient(
		config,
		connection,
	)

	runWithChannels(client)
}

func runWithChannels(client *chat.Client) {
	disconnected := make(chan struct{})
	connected := make(chan struct{})
	message := make(chan string)

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
			case <-connected:
				fmt.Println("Connected")
			case newMessage := <-message:
				fmt.Println(newMessage)
			}
		}
	}()

	client.ConnectWithChannels(connected, disconnected, message)
}

func runWithCallbacks(client *chat.Client) {
	client.ConnectWithCallbacks(
		func() {
			fmt.Println("Connected")
		},
		func() {
			fmt.Println("Disconnected")
		},
		func(message string) {
			fmt.Println(message)
		},
	)
}
