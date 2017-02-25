package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	irc "github.com/fluffle/goirc/client"
	"log"
)

var (
	config configuration.Configuration
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

func initTwitchClient() (client *irc.Conn) {
	client = irc.Client(initTwitchIrcConfig())
	return
}

func init() {
	initConfiguration()
}

func main() {
	client := initTwitchClient()
	quit := make(chan bool)

	client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		client.Join("#" + config.Channel)
	})

	client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})

	client.HandleFunc("privmsg", func(conn *irc.Conn, line *irc.Line) {
		println(line.Args[1])
	})

	for {
		if err := client.Connect(); err != nil {
			fmt.Printf("Connection error: %s\n", err)
			return
		}
		<-quit
	}
}
