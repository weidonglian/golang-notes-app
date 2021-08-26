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
	conn     *nats.Conn
	logger   *logrus.Logger
	ctxClose context.Context
}

func (n *natsSubClient) Close() {
	n.logger.Infof("request draining the nats subscriber client")
	n.conn.Drain()
	n.logger.Infof("wait draining the nats subscriber client to be closed")
	<-n.ctxClose.Done()
	n.logger.Infof("the nats subscriber client is closed now")
}

func (n *natsSubClient) Subscribe(ctx context.Context, key SubjectKey, handler SubscriptionHandler) error {
	if handler == nil {
		panic("pubsub: handler can not be nil")
	}
	_, err := n.conn.Subscribe(key.String(), func(msg *nats.Msg) {
		n.logger.Infof("receive event: [%s] payload: %s", msg.Subject, string(msg.Data))
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
		if err != nil {
			logger.Errorf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
		} else {
			logger.Infof("Disconnected, will attempt reconnects for %.0fm", totalWait.Minutes())
		}
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logger.Infof("Reconnected [%s]", nc.ConnectedUrl())
	}))
	ctxClose, fnClose := context.WithCancel(context.Background())
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		if err := nc.LastError(); err != nil {
			logger.Errorf("exiting nats subscription client with error: %v", err)
		} else {
			logger.Infof("exiting nats subscription client successfully")
		}
		fnClose()
	}))

	conn, err := nats.Connect(cfg.Nats.URL, opts...)

	if err != nil {
		return nil, err
	}

	return &natsSubClient{
		conn:     conn,
		logger:   logger,
		ctxClose: ctxClose,
	}, nil
}
