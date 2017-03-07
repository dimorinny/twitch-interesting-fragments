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
	"os"
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
	timeBuffer := buffer.NewMessagesBuffer(message, time.Second*10)

	go ircChatExample(chat, message)

	bufferedChannel := timeBuffer.Start()

	file, err := os.Create("/Users/damerkurev/Desktop/dump.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for items := range bufferedChannel {

		fmt.Println(
			fmt.Sprintf(
				"%s;%d",
				time.Now().Format("15:04:05"),
				len(items),
			),
		)

		file.WriteString(
			fmt.Sprintf(
				"%s;%d\n",
				time.Now().Format("15:04:05"),
				len(items),
			),
		)
	}
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
