package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-chat-api"
	"github.com/dimorinny/twitch-interesting-fragments/api"
	"github.com/dimorinny/twitch-interesting-fragments/buffer"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	"github.com/dimorinny/twitch-interesting-fragments/detection"
	"log"
	"net/http"
	"time"
)

var (
	config   configuration.Configuration
	twitch   *twitchchat.Chat
	uploader *api.Uploader
)

func initConfiguration() {
	config = configuration.Configuration{}

	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func initChat() {
	twitch = twitchchat.NewChat(
		twitchchat.NewConfiguration(
			config.Nickname,
			config.Oauth,
			config.Channel,
		),
	)
}

func initUploader() {
	uploader = api.NewUploader(config, http.DefaultClient)
}

func init() {
	initConfiguration()
	initChat()
	initUploader()
}

func main() {
	startDetection()
}

func startDetection() {
	messages := make(chan string)
	timeBuffer := buffer.NewMessagesBuffer(
		messages,
		time.Second*time.Duration(config.MessagesBufferTime),
	)

	if err := handleTwitchChat(messages); err != nil {
		log.Fatal(err)
	}

	bufferedChannel := timeBuffer.Start()
	output := make(chan int)

	go func() {
		for item := range output {
			fmt.Printf("Splash detected: %d\n", item)
		}
	}()

	detection.StartDetection(
		config.WindowSize,
		config.SpikeRate,
		config.SmoothRate,
		bufferedChannel,
		output,
	)
}

func handleTwitchChat(message chan string) error {
	disconnected := make(chan struct{})
	connected := make(chan struct{})

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
			case <-connected:
				fmt.Println("Connected")
			}
		}
	}()

	return twitch.ConnectWithChannels(connected, disconnected, message)
}
