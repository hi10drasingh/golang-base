package drmrmq

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Action is an action that occurs after processed this delivery.
type Action int

// Handler defines the handler of each Delivery and return Action.
type Handler func(d amqp.Delivery) (action Action)

const (
	// Ack default ack this msg after you have successfully processed this delivery.
	Ack Action = iota
	// NackDiscard the message will be dropped or delivered to a server configured dead-letter queue.
	NackDiscard
	// NackRequeue deliver this message to a different consumer.
	NackRequeue
)

var (
	errAck         = errors.New("error in ack")
	errNackDiscard = errors.New("error in nack discard")
	errNackRequeue = errors.New("error in nack requeue")
)

// Consumer allows you to create and connect to queues for data consumption.
type Consumer struct {
	rq *RabbitMQ
	ch *amqp.Channel
}

// NewConsumer returns new instance of Consumer.
func (rq *RabbitMQ) NewConsumer(ctx context.Context) (*Consumer, error) {
	channel, err := rq.channel(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Consumer Channel Creation")
	}

	return &Consumer{
		rq: rq,
		ch: channel,
	}, nil
}

// Consume fetches msgs from amqp queue and pass that to Handler.
func (con *Consumer) Consume(ctx context.Context, queueName, tag string, handler Handler) error {
	deliveries, err := con.ch.Consume(
		queueName, // name
		tag,       // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Queue consume")
	}

	go con.processEntries(ctx, deliveries, handler)

	return nil
}

func (con *Consumer) processEntries(ctx context.Context, msgs <-chan amqp.Delivery, handler Handler) {
	for msg := range msgs {
		switch handler(msg) {
		case Ack:
			err := msg.Ack(false)
			if err != nil {
				con.rq.log.Error(ctx, err, errAck.Error())
			}
		case NackDiscard:
			err := msg.Nack(false, false)
			if err != nil {
				con.rq.log.Error(ctx, err, errNackDiscard.Error())
			}
		case NackRequeue:
			err := msg.Nack(false, true)
			if err != nil {
				con.rq.log.Error(ctx, err, errNackRequeue.Error())
			}
		}
	}
}
