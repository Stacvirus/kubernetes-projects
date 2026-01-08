package main

import (
	"broadcast/config"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	config := config.Load()
	// Connect to NATS server
	nc, err := nats.Connect(config.NatsURL, nats.UserInfo(config.NatsUser, config.NatsPassword))
	if err != nil {
		log.Panic("Nats connection failed:", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server")

	// Subscribe to a subject
	subject := config.NatsSubject
	_, err = nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Received message: %s", string(m.Data))
	})
	if err != nil {
		log.Panic("Subscribe to Nats server subject failed:", err)
	}

	log.Printf("Subscribed to %s subject", subject)

	// Keep the connection alive
	select {}
}
