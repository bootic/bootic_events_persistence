package keen

import (
  data "github.com/bootic/bootic_go_data"
  "time"
	"log"
)

type BufferedClient struct {
  projectId      string
  apiKey 				 string
  Notifier       data.EventsChannel
  ticker         *time.Ticker
}

func (cl *BufferedClient) Listen() {
  for {
		select {
		case event := <- cl.Notifier:
	    evtType, _     := event.Get("type").String()
	    log.Println("Event!", evtType)
		case <- cl.ticker.C:
			log.Println("TICK!")
		}
    
  }
}

func NewBufferedClient(projectId, apiKey string, duration time.Duration) (client *BufferedClient, err error) {

  client = &BufferedClient{
    projectId: projectId,
    apiKey: apiKey,
    Notifier: make(data.EventsChannel, 1),
    ticker: time.NewTicker(duration),
  }

  return
}
