package drmrmq

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmtime"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds conf fro amqp server.
type Config struct {
	Host      string             `json:"host" validate:"required"`
	Port      int                `json:"port" validate:"required,number"`
	User      string             `json:"user" validate:"required"`
	Password  string             `json:"password" validate:"required"`
	Vhost     string             `json:"vhost" validate:"required"`
	Timeout   drmtime.CustomTime `json:"timeout" validate:"required"`
	Heartbeat drmtime.CustomTime `json:"heartbeat" validate:"required"`
	Enabled   bool               `json:"enabled"`
}

// RabbitMQ struct contains Connection & Channel.
type RabbitMQ struct {
	conn      *amqp.Connection
	rmqConfig *Config
	log       drmlog.Logger
	mu        *sync.Mutex
}

// NewRabbitMQ return new instance of RabbitMQ Struct.
func NewRabbitMQ(conf *Config, log drmlog.Logger) (*RabbitMQ, error) {
	rmq := RabbitMQ{
		conn:      nil,
		rmqConfig: conf,
		log:       log,
		mu:        &sync.Mutex{},
	}

	if err := rmq.Connect(); err != nil {
		return nil, errors.Wrap(err, "RabbitMQ Conenct")
	}

	log.Debug(context.Background(), "Rabbitmq connected")

	return &rmq, nil
}

// Connect connects to amqp server.
func (rq *RabbitMQ) Connect() (err error) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	connectionTimeout := rq.rmqConfig.Timeout.Time

	address := fmt.Sprintf("amqp://%s:%s@%s:%d",
		rq.rmqConfig.User,
		rq.rmqConfig.Password,
		rq.rmqConfig.Host,
		rq.rmqConfig.Port,
	)

	conn, err := amqp.DialConfig(address, amqp.Config{
		Vhost:     rq.rmqConfig.Vhost,
		Heartbeat: rq.rmqConfig.Heartbeat.Time,
		Dial: func(network, addr string) (conn net.Conn, err error) {
			conn, err = net.DialTimeout(network, addr, connectionTimeout)

			return conn, errors.Wrap(err, "AMQP Connection Timeout")
		},
	})
	if err != nil {
		return errors.Wrap(err, "AMQP Open Connection")
	}

	rq.conn = conn

	// testing channel
	channel, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "AMPQ Channel Creation")
	}

	defer channel.Close()

	go rq.handleBlocking()

	go rq.handleClose()

	return nil
}

// Close closes all open connections.
func (rq *RabbitMQ) Close(ctx context.Context) error {
	if err := rq.conn.Close(); err != nil {
		return errors.Wrap(err, "RabbitMQ Connection Close")
	}

	return nil
}

func (rq *RabbitMQ) handleBlocking() {
	blockings := rq.conn.NotifyBlocked(make(chan amqp.Blocking))

	for b := range blockings {
		if b.Active {
			rq.log.Infof(context.Background(), "AMPQ TCP Conn blocked: %q", b.Reason)
		} else {
			rq.log.Info(context.Background(), "AMPQ TCP Conn unblocked")
		}
	}
}

func (rq *RabbitMQ) handleClose() {
	err := <-rq.conn.NotifyClose(make(chan *amqp.Error))

	if err != nil && err.Server {
		err := rq.Connect()
		if err != nil {
			rq.log.Error(context.Background(), err, "RMQ Connection Reconnect")
		}
	}
}

func (rq *RabbitMQ) channel(ctx context.Context) (*amqp.Channel, error) {
	channel, err := rq.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "AMPQ Channel Creation")
	}

	go rq.closeWithContext(ctx, channel)

	go rq.startNotifyCancelOrClosed(ctx, channel)

	return channel, nil
}

func (rq *RabbitMQ) closeWithContext(ctx context.Context, channel *amqp.Channel) {
	<-ctx.Done()

	channel.Close()
}

func (rq *RabbitMQ) startNotifyCancelOrClosed(ctx context.Context, channel *amqp.Channel) {
	notifyCloseChan := channel.NotifyClose(make(chan *amqp.Error))
	notifyCancelChan := channel.NotifyCancel(make(chan string))

	select {
	case err := <-notifyCloseChan:
		rq.log.Error(ctx, err, "Notify Channel Close")

	case err := <-notifyCancelChan:
		rq.log.Error(ctx, errors.New(err), " Notify Channel Cancel")
	}
}
