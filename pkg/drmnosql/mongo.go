package drmnosql

import (
	"context"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds init dependencies for New DB.
type Config struct {
	MongoConfig config.MongoConfig
	Log         drmlog.Logger
}

// GetDB open and pings connection to provided databases
// and return pointer DB struct with errors.
func GetDB(conf *Config) (*mongo.Client, error) {
	connectionTimeout := conf.MongoConfig.ConnectionTimeout.Time

	clientOpts := options.Client()
	clientOpts = clientOpts.SetHosts(conf.MongoConfig.Hosts).SetAuth(options.Credential{
		Username:   conf.MongoConfig.User,
		Password:   conf.MongoConfig.Password,
		AuthSource: conf.MongoConfig.AuthSource,
	}).SetConnectTimeout(connectionTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, errors.Wrap(err, "NoSQL Database Open Connection")
	}

	conf.Log.Info(context.Background(), "NoSQl Database Connected")

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NoSQL Database Connection Testing")
	}

	conf.Log.Info(context.Background(), "NoSQl Database Tested")

	return client, nil
}
