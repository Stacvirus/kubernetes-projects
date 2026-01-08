package nats

import (
	"errors"
	"log"
	"todo-backend/internal/config"

	"github.com/nats-io/nats.go"
)

type NatClient struct {
	Conn    *nats.Conn
	Subject string
}

func NewNatsClient(config config.Config) (*NatClient, error) {
	nc, err := connect(config)
	if err != nil {
		log.Println("failed to connect to NATS: ", err)
		return nil, err
	}
	return &NatClient{
		Conn:    nc,
		Subject: config.NatsSubject,
	}, nil
}

func connect(config config.Config) (*nats.Conn, error) {
	nc, err := nats.Connect(config.NatsURL, nats.UserInfo(config.NatsUser, config.NatsPassword))
	if err != nil {
		return nil, err
	}
	log.Println("✅ connected to NAT server at: ", config.NatsURL)
	return nc, nil
}

func (c *NatClient) Publish(message string) error {
	if c.Subject == "" {
		return errors.New("nats subject is not set")
	}
	err := c.Conn.Publish(c.Subject, []byte(message))
	if err != nil {
		return err
	}
	log.Printf("Published message to subject %s", c.Subject)
	return nil
}
