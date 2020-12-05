package pubsub

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var ErrPublishEmptySubject = fmt.Errorf("won't publish an empty subject key")

type PubClient interface {
	Publish(ctx context.Context, key SubjectKey, data interface{}) error
}

var _ PubClient = &natsPubClient{}

type natsPubClient struct {
	conn   *nats.Conn
	logger *logrus.Logger
}

func (n *natsPubClient) Publish(ctx context.Context, key SubjectKey, data interface{}) error {
	if key.IsEmpty() {
		return ErrPublishEmptySubject
	}
	return n.conn.Publish(key.String(), nil)
}

func NewPubClient(logger *logrus.Logger, natsConn *NatsConnection) (PubClient, error) {
	return &natsPubClient{
		conn:   natsConn.Conn,
		logger: logger,
	}, nil
}
