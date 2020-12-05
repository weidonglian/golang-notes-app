package pubsub

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var ErrSubscribeEmptySubject = fmt.Errorf("won't subscribe an empty subject key")

type SubClient interface {
	Subscribe(ctx context.Context, key SubjectKey, data interface{}) error
}

var _ SubClient = &natsSubClient{}

type natsSubClient struct {
	conn   *nats.Conn
	logger *logrus.Logger
}

func (n *natsSubClient) Subscribe(ctx context.Context, key SubjectKey, data interface{}) error {
	panic("implement me")
}

func NewSubClient(logger *logrus.Logger, natsConn *NatsConnection) (SubClient, error) {
	return &natsSubClient{
		conn:   natsConn.Conn,
		logger: logger,
	}, nil
}
