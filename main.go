package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dimorinny/twitch-chat-api"
	"github.com/dimorinny/twitch-interesting-fragments/api"
	"github.com/dimorinny/twitch-interesting-fragments/buffer"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	"github.com/dimorinny/twitch-interesting-fragments/detection"
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
	bufferExample()
}

func detectionExample() {
	input := make(chan int)
	output := make(chan int)

	data := []int{
		2, 7, 8, 6, 4, 4, 5, 7, 9, 6, 7, 9, 6, 7, 5, 5, 6, 5, 4, 8, 5, 7, 10,
		9, 5, 7, 8, 15, 11, 11, 5, 4, 8, 4, 8, 6, 4, 13, 12, 10, 5, 10, 11, 12,
		11, 6, 9, 7, 9, 7, 29, 18, 17, 17, 9, 10, 8, 14, 8, 10, 10, 13, 13, 10, 10,
	}

	go func() {
		for _, item := range data {
			input <- item
		}
		close(input)
	}()

	go func() {
		for item := range output {
			fmt.Printf("Splash detected: %d\n", item)
		}
	}()

	detection.StartDetection(
		config.WindowSize,
		config.SpikeRate,
		config.SmoothRate,
		input,
		output,
	)
}

func uploadExample() {
	result, err := uploader.Upload()

	if err != nil {
		log.Fatal(err)
	}

	print(result.Url)
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
	timeBuffer := buffer.NewMessagesBuffer(message, time.Second*25)

	go ircChatExample(chat, message)

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

func ircChatExample(twitch *twitchchat.Chat, message chan string) {
	stop := make(chan struct{})
	defer close(stop)

	disconnected := make(chan struct{})
	connected := make(chan struct{})

	go func() {
		for {
			select {
			case <-disconnected:
				fmt.Println("Disconnected")
				stop <- struct{}{}
			case <-connected:
				fmt.Println("Connected")
			}
		}
	}()

	err := twitch.ConnectWithChannels(connected, disconnected, message)
	if err != nil {
		return
	}

	<-stop
}
