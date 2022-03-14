package drmnosql

import (
	"context"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmtime"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds connection details of mongo server.
type Config struct {
	// AuthSource        string     `json:"authSource" validate:"required"`
	Hosts             []string           `json:"hosts" validate:"required"`
	User              string             `json:"user" validate:"required"`
	Password          string             `json:"password" validate:"required"`
	DB                string             `json:"db" validate:"required"`
	ConnectionTimeout drmtime.CustomTime `json:"connectionTimeout" validate:"required"`
}

// GetDB open and pings connection to provided databases
// and return pointer DB struct with errors.
func GetDB(conf *Config, log drmlog.Logger) (*mongo.Client, error) {
	connectionTimeout := conf.ConnectionTimeout.Time

	clientOpts := options.Client()
	clientOpts = clientOpts.SetHosts(conf.Hosts).SetAuth(options.Credential{
		Username: conf.User,
		Password: conf.Password,
		// AuthSource: conf.AuthSource,
	}).SetConnectTimeout(connectionTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, errors.Wrap(err, "NoSQL Database Open Connection")
	}

	log.Debug(context.Background(), "NoSQl Database Connected")

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NoSQL Database Connection Testing")
	}

	log.Debug(context.Background(), "NoSQl Database Tested")

	return client, nil
}
