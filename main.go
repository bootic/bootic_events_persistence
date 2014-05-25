package main

import (
	"bootic_events_persistence/backend"
	"bootic_events_persistence/keenio"
	"flag"
	data "github.com/bootic/bootic_go_data"
	bootic_sse "github.com/bootic/bootic_sse"
	bootic_zmq "github.com/bootic/bootic_zmq"
	"log"
	"time"
)

type BooticStream interface {
	SubscribeToType(data.EventsChannel, string)
}

func main() {
	var (
		interval      string
		keenProjectId string
		keenApiKey    string
		httpToken     string
		httpUrl       string
		transport     string
		topic         string
		zmqAddress    string
	)

	flag.StringVar(&interval, "interval", "60s", "Time interval to send stats on. Ie. 30s, 2m, etc")
	flag.StringVar(&keenProjectId, "keenprojectid", "", "Keen.io project id")
	flag.StringVar(&keenApiKey, "keenapikey", "", "Keen.io API Key")
	flag.StringVar(&httpToken, "httptoken", "", "Token to use if transport is httpstream")
	flag.StringVar(&httpUrl, "httpurl", "", "URL to use if transport is httpstream")
	flag.StringVar(&transport, "transport", "zmq", "What stream transport to listen to (zmq or httpstream)")
	flag.StringVar(&topic, "topic", "all", "Stream topic to subscribe to") // event type. ie "order", "pageview"
	flag.StringVar(&zmqAddress, "zmqsocket", "tcp://127.0.0.1:6000", "ZMQ socket address to bind to")

	flag.Parse()

	duration, err := time.ParseDuration(interval)
	if err != nil {
		panic("INTERVAL cannot be parsed")
	}

	var stream BooticStream

	// Setup ZMQ or HTTP subscriber +++++++++++++++++++++++++++++++
	switch transport {
	case "httpstream":
		stream, err = bootic_sse.NewClient(httpUrl, httpToken)
		log.Println("Listening on HTTP stream", httpUrl)
	default:
		stream, err = bootic_zmq.NewZMQSubscriber(zmqAddress, "")
		log.Println("Listening on ZMQ stream", zmqAddress)
	}

	if err != nil {
		panic(err)
	}

	// multi-backend buffered client
	backends, err := backend.NewBufferedClient(duration)
	if err != nil {
		panic(err)
	}

	// Keen.io backend
	keenCl := keenio.NewClient(keenProjectId, keenApiKey)
	backends.Register(keenCl)

	// Push events to all registered backends
	stream.SubscribeToType(backends.Notifier, topic)

	log.Println("Sending events to Keenio as", keenProjectId)
	backends.Listen()

}
