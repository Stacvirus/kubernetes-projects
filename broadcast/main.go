package main

import (
	"broadcast/config"
	"broadcast/telegram"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	config := config.Load()

	// Telegram client
	tg := telegram.New(
		config.TelegramBotToken,
		config.TelegramChatID,
	)

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
		msg := string(m.Data)

		log.Printf("Received message: %s", msg)

		err := tg.SendMessage(msg)
		if err != nil {
			log.Println("❌ Telegram send failed:", err)
		} else {
			log.Println("📤 Message forwarded to Telegram")
		}
	})

	if err != nil {
		log.Panic("Subscribe to Nats server subject failed:", err)
	}

	log.Printf("Subscribed to %s subject", subject)

	// Keep the connection alive
	select {}
}
