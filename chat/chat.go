package chat

import (
	"fmt"
	"github.com/dimorinny/twitch-interesting-fragments/configuration"
	irc "github.com/fluffle/goirc/client"
)

const (
	connectedEvent    = "connected"
	disconnectedEvent = "disconnected"
	newMessageEvent   = "privmsg"
)

type (
	Connected    func()
	Disconnected func()
	NewMessage   func(message string)

	Client struct {
		config     configuration.Configuration
		connection *irc.Conn
	}
)

func NewClient(config configuration.Configuration, connection *irc.Conn) *Client {
	return &Client{
		config:     config,
		connection: connection,
	}
}

func (c *Client) ConnectWithChannels(connected, disconnected chan<- struct{}, message chan<- string) {
	connectedCallback := func() {
		connected <- struct{}{}
	}

	disconnectedCallback := func() {
		disconnected <- struct{}{}
	}

	newMessageCallback := func(newMessage string) {
		message <- newMessage
	}

	c.ConnectWithCallbacks(connectedCallback, disconnectedCallback, newMessageCallback)
}

func (c *Client) ConnectWithCallbacks(connected Connected, disconnected Disconnected, message NewMessage) {
	quit := make(chan struct{})

	c.connection.HandleFunc(connectedEvent, func(conn *irc.Conn, line *irc.Line) {
		connected()
		c.connection.Join("#" + c.config.Channel)
	})

	c.connection.HandleFunc(disconnectedEvent, func(conn *irc.Conn, line *irc.Line) {
		disconnected()
		quit <- struct{}{}
	})

	c.connection.HandleFunc(newMessageEvent, func(conn *irc.Conn, line *irc.Line) {
		if len(line.Args) > 1 {
			message(line.Args[1])
		}
	})

	if err := c.connection.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err)
		disconnected()
		return
	}

	<-quit
}
