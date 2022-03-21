package drmrmq

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds config for publish.
type PublisherConfig struct {
	Exchange string
	Key      string
	Msg      amqp.Publishing
}

type Publisher struct {
	rq *RabbitMQ
	ch *amqp.Channel
}

func (rq *RabbitMQ) NewPublisher(ctx context.Context) (*Publisher, error) {
	pub := &Publisher{
		rq: rq,
	}

	if err := pub.channel(ctx); err != nil {
		return nil, errors.Wrap(err, "Publisher Channel Creation")
	}

	return pub, nil
}

func (pub *Publisher) channel(ctx context.Context) error {
	channel, err := pub.rq.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "AMPQ Channel Creation")
	}

	pub.ch = channel

	go pub.closeWithContext(ctx)

	go pub.startNotifyCancelOrClosed(ctx)

	return nil
}

// Publish will send the provided msg to provided exchange.
// It also acknowledge the publish.
func (pub *Publisher) Publish(ctx context.Context, config *PublisherConfig) (err error) {
	// Confirming if published
	if err = pub.ch.Confirm(false); err != nil {
		return errors.Wrap(err, "Queue Publish")
	}

	confirms := pub.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	defer pub.confirmOne(ctx, confirms)

	err = pub.ch.Publish(
		config.Exchange, // exchange
		config.Key,      // routing key
		false,           // mandatory
		false,           // immediate
		config.Msg,
	)

	return errors.Wrap(err, "Queue Publish")
}

func (pub *Publisher) confirmOne(ctx context.Context, confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		pub.rq.log.Infof(ctx, "confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		pub.rq.log.Infof(ctx, "failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

func (pub *Publisher) closeWithContext(ctx context.Context) {
	<-ctx.Done()

	if err := pub.ch.Close(); err != nil {
		pub.rq.log.Error(ctx, err, "Publisher Close")
	}
}

func (pub *Publisher) startNotifyCancelOrClosed(ctx context.Context) {
	notifyCloseChan := pub.ch.NotifyClose(make(chan *amqp.Error))
	notifyCancelChan := pub.ch.NotifyCancel(make(chan string))

	select {
	case err := <-notifyCloseChan:
		if err != nil && err.Server {
			er := pub.channel(ctx)
			if er != nil {
				pub.rq.log.Error(ctx, er, "Publisher Channel Reconnect")
			}

			pub.rq.log.Error(ctx, err, "Publisher Notify Channel Close")
		}

	case err := <-notifyCancelChan:
		er := pub.channel(ctx)
		if er != nil {
			pub.rq.log.Error(ctx, er, "Publisher Channel Reconnect")
		}

		pub.rq.log.Error(ctx, errors.New(err), "Publisher Notify Channel Close")
	}
}
