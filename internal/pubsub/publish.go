package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
)

var ErrPublishEmptySubject = fmt.Errorf("won't publish an empty subject key")

type Publisher interface {
	Publish(ctx context.Context, key SubjectKey, data interface{}) error
	Close()
}

var _ Publisher = &natsPubClient{}

type natsPubClient struct {
	conn   *nats.Conn
	logger *logrus.Logger
}

func (n *natsPubClient) Close() {
	n.conn.Close()
}

func (n *natsPubClient) Publish(ctx context.Context, key SubjectKey, data interface{}) error {
	if key.IsEmpty() {
		return ErrPublishEmptySubject
	}

	var (
		b   []byte
		err error
	)

	b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	err = n.conn.Publish(key.String(), b)

	if err != nil {
		n.logger.Error(fmt.Errorf("couldn't publish to nats: %s", err.Error()))
	}

	return err
}

func NewPubClient(logger *logrus.Logger, cfg *config.Config) (Publisher, error) {
	opts := []nats.Option{nats.Name("Notes-App NATS Publisher")}
	conn, err := nats.Connect(cfg.Nats.URL, opts...)

	if err != nil {
		return nil, err
	}

	return &natsPubClient{
		conn:   conn,
		logger: logger,
	}, nil
}
