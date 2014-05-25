package keen

import (
  data "github.com/bootic/bootic_go_data"
  "time"
	"log"
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

type StatsClient interface {
	AddEvent(*data.Event)
	Submit()
}

const(
	KeenUrl = "https://api.keen.io/3.0/projects/%s/events?api_key=%s"
)

type Client struct {
  projectId			string
  apiKey				string
	url						string
	buffer				map[string][]interface{}
	eventCount		int
}

type metaData struct {
	Timestamp string `json:"timestamp"`
}

func (c *Client) AddEvent(event *data.Event) {
	evtType, _  := event.Get("type").String()
	evtTime, _  := event.Get("time").String()
	evtPayload, _ := event.Get("data").Map()

	evtPayload["keen"] = &metaData{evtTime}
	c.buffer[evtType] = append(c.buffer[evtType], evtPayload)
	c.eventCount = c.eventCount + 1
}

func (c *Client) Submit() {
	if c.eventCount == 0 {
		return
	}

	buff := c.buffer
	evtCount := c.eventCount
	c.reset()
	jsonBytes, err := json.Marshal(buff)
	if err != nil {
		log.Fatal("Could not encode JSON", err, buff)
	}

	resp, err := http.Post(c.url, "application/json", bytes.NewReader(jsonBytes))
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Could not POST JSON", err)
	}

	log.Println(fmt.Sprintf("Sent %v events to KeenIo. Resp %v.", evtCount, resp.StatusCode))
}

func (c *Client) reset() {
	c.buffer = make(map[string][]interface{})
	c.eventCount = 0
}

func NewClient(projectId, apiKey string) (client *Client) {

	client = &Client{
		projectId: projectId,
		apiKey: apiKey,
		url: fmt.Sprintf(KeenUrl, projectId, apiKey),
	}

	client.reset()

	return
}

type BufferedClient struct {
  Notifier			data.EventsChannel
  ticker				*time.Ticker
	clients				[]StatsClient
}

func (cl *BufferedClient) Listen() {
  for {
		select {
		case event := <- cl.Notifier:
			for i := range cl.clients {
				cl.clients[i].AddEvent(event)
			}
		case <- cl.ticker.C:
			for i := range cl.clients {
				go cl.clients[i].Submit()
			}
		}
    
  }
}

func (cl *BufferedClient) Register(c StatsClient) {
	cl.clients = append(cl.clients, c)
}

func NewBufferedClient(projectId, apiKey string, duration time.Duration) (client *BufferedClient, err error) {

  client = &BufferedClient{
    Notifier: make(data.EventsChannel, 1),
    ticker: time.NewTicker(duration),
  }

	client.Register(NewClient(projectId, apiKey))

  return
}
