package main

import (
  "bootic_stathat/client"
  "flag"
  "github.com/bootic/bootic_zmq"
  "log"
  "time"
)

func main() {
  var (
    topic          string
    zmqAddress     string
    stathatAccount string
    interval       string
  )

  flag.StringVar(&topic, "topic", "", "ZMQ topic to subscribe to") // event type. ie "order", "pageview"
  flag.StringVar(&zmqAddress, "zmqsocket", "tcp://127.0.0.1:6000", "ZMQ socket address to bind to")
  flag.StringVar(&stathatAccount, "stathatAccount", "", "Stathat email or account key")
  flag.StringVar(&interval, "interval", "60s", "Time interval to send stats on. Ie. 30s, 2m, etc")

  flag.Parse()

  duration, err := time.ParseDuration(interval)
  if err != nil {
    panic("INTERVAL cannot be parsed")
  }

  // Setup ZMQ subscriber +++++++++++++++++++++++++++++++
  zmq, _ := booticzmq.NewZMQSubscriber(zmqAddress, topic)

  log.Println("ZMQ socket started on", zmqAddress, "topic '", topic, "'")

  cl, err := client.NewBufferedClient(stathatAccount, topic, duration)
  if err != nil {
    panic("Client could not connect")
  }

  log.Println("Sending", topic, "events to Stathat as", stathatAccount)
  zmq.SubscribeToType(cl.Notifier, topic)

  cl.Listen()

}
