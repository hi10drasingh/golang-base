package drmrmq

import (
	"fmt"
	"net"
	"time"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ struct contains Connection & Channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Chan       *amqp.Channel
	RMQConfig  config.RabbitMQConfig
	Log        drmlog.Logger
}

// NewRabbitMQ return new instance of RabbitMQ Struct
func NewRabbitMQ(conf config.RabbitMQConfig, log drmlog.Logger) (*RabbitMQ, error) {
	rq := RabbitMQ{}

	rq.RMQConfig = conf

	rq.Log = log

	err := rq.Connect()

	return &rq, err

}

// Connect connects to amqp server
func (rq *RabbitMQ) Connect() (err error) {
	connectionTimeout := time.Duration(rq.RMQConfig.Timeout)

	address := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		rq.RMQConfig.User,
		rq.RMQConfig.Password,
		rq.RMQConfig.Host,
		rq.RMQConfig.Port,
	)

	conn, err := amqp.DialConfig(address, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connectionTimeout)
		},
	})

	if err != nil {
		return errors.Wrap(err, "AMQP Open Connection")
	}

	rq.Connection = conn

	ch, err := rq.Connection.Channel()

	if err != nil {
		return errors.Wrap(err, "AMQP Open Channel")
	}

	rq.Chan = ch

	return
}

// Close closes all open connections
func (rq *RabbitMQ) Close() error {
	err := rq.Chan.Close()

	if err != nil {
		return errors.Wrap(err, "RabbitMQ Channel Close")
	}

	err = rq.Connection.Close()

	if err != nil {
		return errors.Wrap(err, "RabbitMQ Connection Close")
	}

	return nil
}
