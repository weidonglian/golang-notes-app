package pubsub

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
	"time"
)

var ErrSubscribeEmptySubject = fmt.Errorf("won't subscribe an empty subject key")

type Subscriber interface {
	Subscribe(ctx context.Context, key SubjectKey, handler SubscriptionHandler) error
	Close()
}

type SubscriptionHandler func(msg *nats.Msg)

var _ Subscriber = &natsSubClient{}

type natsSubClient struct {
	conn   *nats.Conn
	logger *logrus.Logger
}

func (n *natsSubClient) Close() {
	n.conn.Drain()
}

func (n *natsSubClient) Subscribe(ctx context.Context, key SubjectKey, handler SubscriptionHandler) error {
	if handler == nil {
		panic("pubsub: handler can not be nil")
	}
	_, err := n.conn.Subscribe(key.String(), func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return err
	}
	return n.conn.Flush()
}

func NewSubClient(logger *logrus.Logger, cfg *config.Config) (Subscriber, error) {
	opts := []nats.Option{nats.Name("Notes-App Nats Subscriber")}

	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logger.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logger.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logger.Fatalf("Exiting: %v", nc.LastError())
	}))

	conn, err := nats.Connect(cfg.Nats.URL, opts...)

	if err != nil {
		return nil, err
	}

	return &natsSubClient{
		conn:   conn,
		logger: logger,
	}, nil
}
