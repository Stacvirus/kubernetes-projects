package main

import (
	"broadcast/config"
	"broadcast/telegram"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

func main() {
	config := config.Load()

	var natsHealthy atomic.Bool

	// Telegram client
	tg := telegram.New(
		config.TelegramBotToken,
		config.TelegramChatID,
	)

	// Connect to NATS server
	nc, err := nats.Connect(
		config.NatsURL,
		nats.UserInfo(config.NatsUser, config.NatsPassword),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Println("❌ Disconnected from NATS server:", err)
			natsHealthy.Store(false)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Println("✅ Reconnected to NATS server")
			natsHealthy.Store(true)
		}),
	)
	if err != nil {
		log.Panic("Nats connection failed:", err)
	}
	defer nc.Close()

	natsHealthy.Store(true)
	log.Println("✅ Connected to NATS server")

	// Health endpoint
	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			if !natsHealthy.Load() {
				log.Println("❌ Health check failed: NATS unhealthy")
				http.Error(w, "NATS unhealthy", http.StatusServiceUnavailable)
				return
			}
			log.Println("🩺 Health check OK")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		log.Printf("🩺 Health endpoint listening on :%s/healthz", config.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
	}()

	// Subscribe to a subject
	queueGroup := "broadcaster-workers"
	subject := config.NatsSubject
	_, err = nc.QueueSubscribe(subject, queueGroup, func(m *nats.Msg) {
		msg := string(m.Data)

		log.Printf("Received message: %s", msg)

		if config.Env != "staging" {
			err := tg.SendMessage(msg)
			if err != nil {
				log.Println("❌ Telegram send failed:", err)
			} else {
				log.Println("📤 Message forwarded to Telegram")
			}
		}
	})

	if err != nil {
		log.Panic("Subscribe to Nats server subject failed:", err)
	}

	log.Printf("Subscribed to %s subject", subject)

	// Keep the connection alive
	select {}
}
