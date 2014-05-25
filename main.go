package main

import (
  "flag"
  datasource "github.com/bootic/bootic_sse_client"
  "log"
  "time"
	"bootic_keenio/keen"
)

func main() {
  var (
    interval       string
		projectId			 string
		apiKey			 	 string
  )

  flag.StringVar(&interval, "interval", "60s", "Time interval to send stats on. Ie. 30s, 2m, etc")
  flag.StringVar(&projectId, "projectid", "", "Keen.io project id")
  flag.StringVar(&apiKey, "apikey", "", "Keen.io API Key")

  flag.Parse()

	duration, err := time.ParseDuration(interval)
	if err != nil {
	  panic("INTERVAL cannot be parsed")
	}

  // Setup ZMQ subscriber +++++++++++++++++++++++++++++++
  stream, _ := datasource.NewClient("https://tracker.bootic.net/stream?raw=1", "b00t1csse")

	// Keen.io buffered client
	keenClient, err := keen.NewBufferedClient(projectId, apiKey, duration)
	if err != nil {
	  panic(err)
	}

	stream.Subscribe(keenClient.Notifier)

	log.Println("Sending events to Keenio as", projectId)
	keenClient.Listen()

  // cl, err := client.NewBufferedClient(stathatAccount, topic, duration)
 //  if err != nil {
 //    panic("Client could not connect")
 //  }
 // 
 //  log.Println("Sending", topic, "events to Stathat as", stathatAccount)
 //  zmq.SubscribeToType(cl.Notifier, topic)
 // 
 //  cl.Listen()

}
