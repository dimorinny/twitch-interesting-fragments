package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-chat-api"
	"github.com/dimorinny/twitch-interesting-fragments/api"
	"github.com/dimorinny/twitch-interesting-fragments/buffer"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	"github.com/dimorinny/twitch-interesting-fragments/data"
	"github.com/dimorinny/twitch-interesting-fragments/detection"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
)

var (
	config   configuration.Configuration
	twitch   *twitchchat.Chat
	uploader *api.Uploader
	storage  data.Storage
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

func initMongoStorage() {
	session, err := mgo.Dial(config.StorageHost)
	if err != nil {
		log.Fatal(err)
	}

	storage = data.NewMongoStorage(
		session,
	)
}

func initUploader() {
	uploader = api.NewUploader(config, http.DefaultClient)
}

func init() {
	initConfiguration()
	initChat()
	initMongoStorage()
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

	go handleDetectedFragment(output)

	detection.StartDetection(
		config.WindowSize,
		config.SpikeRate,
		config.SmoothRate,
		bufferedChannel,
		output,
	)
}

func handleDetectedFragment(detectedFragmentChannel <-chan int) {
	for range detectedFragmentChannel {
		response, err := uploader.Upload()
		if err != nil {
			fmt.Printf("Error during uploading fragment %s\n", err)
		}

		storage.AddUploadedFragment(
			data.UploadedFragment{
				Name: response.Name,
				Url:  response.Url,
			},
		)
	}
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
