package pubsub

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var ErrPublishEmptySubject = fmt.Errorf("won't publish an empty subject key")

type PubClient interface {
	Publish(ctx context.Context, key SubjectKey, data proto.Message) error
}

var _ PubClient = &natsPubClient{}

type natsPubClient struct {
	conn   *nats.Conn
	logger *logrus.Logger
}

func (n *natsPubClient) Publish(ctx context.Context, key SubjectKey, data proto.Message) error {
	if key.IsEmpty() {
		return ErrPublishEmptySubject
	}

	var (
		b   []byte
		err error
	)

	b, err = proto.Marshal(data)
	if err != nil {
		return err
	}

	w := &PubSubMessage{
		Data:        b,
		Metadata:    nil,
		PublishTime: ptypes.TimestampNow(),
	}

	b, err = proto.Marshal(w)
	if err != nil {
		return err
	}

	err = n.conn.Publish(key.String(), b)

	if err != nil {
		n.logger.Error(fmt.Errorf("couldn't publish to nats: %s", err.Error()))
	}

	return err
}

func NewPubClient(logger *logrus.Logger, natsConn *NatsConnection) (PubClient, error) {
	return &natsPubClient{
		conn:   natsConn.Conn,
		logger: logger,
	}, nil
}
