package backend

import (
	data "github.com/bootic/bootic_go_data"
	"log"
	"time"
)

type EventsPoster interface {
	AddEvent(*data.Event)
	Submit()
}

type BufferedClient struct {
	Notifier data.EventsChannel
	ticker   *time.Ticker
	clients  []EventsPoster
}

func (cl *BufferedClient) Listen() {
	for {
		select {
		case event := <-cl.Notifier:
			for i := range cl.clients {
				cl.clients[i].AddEvent(event)
			}
		case <-cl.ticker.C:
			for i := range cl.clients {

				go func() {
					defer func() {
						if err := recover(); err != nil {
							log.Println("Goroutine failed:", cl.clients[i], err)
						}
					}()

					cl.clients[i].Submit()
				}()

			}
		}

	}
}

func (cl *BufferedClient) Register(c EventsPoster) {
	cl.clients = append(cl.clients, c)
}

func NewBufferedClient(duration time.Duration) (client *BufferedClient, err error) {

	client = &BufferedClient{
		Notifier: make(data.EventsChannel, 1),
		ticker:   time.NewTicker(duration),
	}

	return
}
