package main

import (
	"context"
	"flag"
	"time"

	"github.com/droomlab/drm-coupon/cmd/consumer/commands"
	"github.com/droomlab/drm-coupon/internal/app/dependency"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
)

type consumeConfig struct {
	queueName   string
	handlerName string
	tagName     string
	lifetime    time.Duration
}

const (
	defaultQueue    = "drmrmq_testing"
	defaultHandler  = "HandleDRMTesting"
	defaultTag      = "simple-consumer"
	defaultLifeTime = 5 * time.Second
)

func main() {
	var logger drmlog.Logger = drmlog.NewConsoleLogger()

	ctx := logger.GetLogger().WithContext(context.Background())

	if err := run(ctx, logger); err != nil {
		logger.Fatal(ctx, err, "Main Error")
	}
}

func run(ctx context.Context, log drmlog.Logger) (err error) {
	consumeConf := parseFlags()

	dependencies, err := dependency.Init()
	if err != nil {
		return errors.Wrap(err, "Initializing dependencies")
	}

	defer func() {
		if er := dependencies.Close(); er != nil {
			err = er
		}
	}()

	err = consumeEntries(ctx, dependencies, consumeConf)

	if err != nil {
		return errors.Wrap(err, "Consumer Entries")
	}

	if consumeConf.lifetime > 0 {
		log.Infof(ctx, "running for %s", consumeConf.lifetime)
		time.Sleep(consumeConf.lifetime)
	} else {
		log.Info(ctx, "running forever")
		select {}
	}

	return err
}

func parseFlags() consumeConfig {
	var (
		queue    = flag.String("queue", defaultQueue, "Ephemeral AMQP queue name")
		handler  = flag.String("handler", defaultHandler, "Msg Delivery Handler")
		tag      = flag.String("consumer-tag", defaultTag, "AMQP consumer tag (should not be blank)")
		lifetime = flag.Duration("lifetime", defaultLifeTime, "lifetime of process before shutdown (0s=infinite)")
	)

	flag.Parse()

	return consumeConfig{
		queueName:   *queue,
		handlerName: *handler,
		tagName:     *tag,
		lifetime:    *lifetime,
	}
}

func consumeEntries(ctx context.Context, dependencies *dependency.Dependency, consumeConf consumeConfig) error {
	consumer, err := dependencies.RMQ.NewConsumer(ctx)
	if err != nil {
		return errors.Wrap(err, "Initializing Consumer")
	}

	comds := commands.GetCommands()

	handlerFunc := comds[consumeConf.handlerName](dependencies)

	err = consumer.Consume(ctx, consumeConf.queueName, consumeConf.tagName, handlerFunc)

	if err != nil {
		return errors.Wrap(err, "Consuming entry")
	}

	return nil
}
