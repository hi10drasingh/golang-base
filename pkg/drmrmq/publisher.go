package drmrmq

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishConfig holds config for publish.
type PublishConfig struct {
	Exchange string
	Key      string
	Msg      amqp.Publishing
}

// Publish will send the provided msg to provided exchange.
// It also acknowledge the publish.
func (rq *RabbitMQ) Publish(ctx context.Context, publishConfig *PublishConfig) (err error) {
	err = rq.CheckEnabled()

	if err != nil {
		return errors.Wrap(err, "Queue Publish")
	}

	err = rq.Connect()

	if err != nil {
		return errors.Wrap(err, "Queue Publish")
	}

	defer func() {
		if er := rq.Close(); er != nil {
			err = errors.Wrap(er, "Queue Publish")
		}
	}()

	if err != nil {
		return errors.Wrap(err, "Queue Publish")
	}

	// Confirming if published
	if err = rq.Chan.Confirm(false); err != nil {
		return errors.Wrap(err, "Queue Publish")
	}

	confirms := rq.Chan.NotifyPublish(make(chan amqp.Confirmation, 1))

	defer rq.confirmOne(ctx, confirms)

	err = rq.Chan.Publish(
		publishConfig.Exchange, // exchange
		publishConfig.Key,      // routing key
		false,                  // mandatory
		false,                  // immediate
		publishConfig.Msg,
	)

	return errors.Wrap(err, "Queue Publish")
}

func (rq *RabbitMQ) confirmOne(ctx context.Context, confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		rq.Log.Infof(ctx, "confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		rq.Log.Infof(ctx, "failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
