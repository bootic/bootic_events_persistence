package client

import (
  data "github.com/bootic/bootic_go_data"
  stathat "github.com/stathat/go"
  "time"
)

type BufferedClient struct {
  eventName      string
  stathatAccount string
  Notifier       data.EventsChannel
  ticker         *time.Ticker
  count          int
}

func (self *BufferedClient) Listen() {
  for {
    <-self.Notifier
    // evtType, _     := event.Get("type").String()
    self.increment()
    // log.Println("Event!", evtType)
  }
}

func (self *BufferedClient) increment() {
  self.count = self.count + 1
}

func (self *BufferedClient) reset() {
  self.count = 0
}

func (self *BufferedClient) post(count int) {
  stathat.PostEZCount(self.eventName, self.stathatAccount, count)
}

func (self *BufferedClient) tick() {
  for {
    select {
    case <-self.ticker.C:
      if self.count != 0 {
        go self.post(self.count)
        self.reset()
      }
    }
  }
}

func NewBufferedClient(stathatAccount string, eventName string, duration time.Duration) (client *BufferedClient, err error) {

  client = &BufferedClient{
    stathatAccount: stathatAccount,
    eventName:      eventName,
    Notifier:       make(data.EventsChannel, 1),
    ticker:         time.NewTicker(duration),
  }

  go client.tick()

  return
}
