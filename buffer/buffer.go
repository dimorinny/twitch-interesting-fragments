package buffer

import (
	"sync"
	"time"
)

type (
	MessagesBuffer struct {
		sync.RWMutex

		ticker   *time.Ticker
		interval time.Duration

		source <-chan string
		items  []string

		inputStop  chan struct{}
		outputStop chan struct{}
	}
)

func NewMessagesBuffer(source <-chan string, interval time.Duration) *MessagesBuffer {
	return &MessagesBuffer{
		ticker:     nil,
		interval:   interval,
		source:     source,
		items:      []string{},
		inputStop:  make(chan struct{}),
		outputStop: make(chan struct{}),
	}
}

func (b *MessagesBuffer) Stop() {
	if b.ticker != nil {
		b.inputStop <- struct{}{}
		b.outputStop <- struct{}{}

		b.ticker.Stop()
		b.ticker = nil
	}
}

func (b *MessagesBuffer) Start() <-chan int {
	b.Stop()

	b.ticker = time.NewTicker(b.interval)
	outputChannel := make(chan int)

	go b.startInputChannelHandling()
	go b.startOutputChannelHandling(outputChannel)

	return outputChannel
}

func (b *MessagesBuffer) startOutputChannelHandling(output chan<- int) {
	for {
		select {
		case _ = <-b.ticker.C:
			b.Lock()
			output <- len(b.items)
			b.items = []string{}
			b.Unlock()

		case _ = <-b.outputStop:
			b.Lock()
			b.items = []string{}
			b.Unlock()

			close(output)
			return
		}
	}
}

func (b *MessagesBuffer) startInputChannelHandling() {
	for {
		select {
		case item, ok := <-b.source:
			if !ok {
				go b.Stop()
			} else {
				b.Lock()
				b.items = append(b.items, item)
				b.Unlock()
			}

		case _ = <-b.inputStop:
			return
		}
	}
}
