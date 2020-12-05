package pubsub

import (
	"github.com/nats-io/nats.go"
	"github.com/weidonglian/notes-app/config"
)

type NatsConnection struct {
	Conn *nats.Conn
}

func NewNatsConnection(cfg *config.Config) (*NatsConnection, error) {
	opts := []nats.Option{nats.Name("NATS Notes-App")}
	conn, err := nats.Connect(cfg.Nats.URL, opts...)

	if err != nil {
		return nil, err
	}

	return &NatsConnection{Conn: conn}, nil
}

func (n *NatsConnection) Close() {
	n.Conn.Close()
}
