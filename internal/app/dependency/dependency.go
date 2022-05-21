package dependency

import (
	"context"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmnosql"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/droomlab/drm-coupon/pkg/drmsql"
	"github.com/pkg/errors"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Server struct holds all server level dependencies.
type Dependency struct {
	Config *config.App
	Log    drmlog.Logger
	SQL    *nap.DB
	NoSQL  *mongo.Client
	RMQ    *drmrmq.RabbitMQ
}

// Init initialized global dependencies.
func Init() (*Dependency, error) {
	conf, err := config.Load()
	if err != nil {
		return nil, errors.Wrap(err, "Config Initialize")
	}

	log, err := drmlog.NewZeroLogger(conf.Log)
	if err != nil {
		return nil, errors.Wrap(err, "Log Initialize")
	}

	sqldb, err := drmsql.GetDB(&conf.Mysql, log)
	if err != nil {
		return nil, errors.Wrap(err, "SQL DB Initialize")
	}

	nosqldb, err := drmnosql.GetDB(&conf.Mongo, log)
	if err != nil {
		return nil, errors.Wrap(err, "NoSQL DB Initialize")
	}

	rmqlog, err := drmlog.NewRMQLogger(conf.Log)
	if err != nil {
		return nil, errors.Wrap(err, "RMQLog Initialize")
	}

	rmq, err := drmrmq.NewRabbitMQ(&conf.RabbitMQ, rmqlog)
	if err != nil {
		return nil, errors.Wrap(err, "RabbitMQ Initialize")
	}

	return &Dependency{
		Config: conf,
		Log:    log,
		SQL:    sqldb,
		NoSQL:  nosqldb,
		RMQ:    rmq,
	}, nil
}

// Close closes all global dependencies.
func (d *Dependency) Close() error {
	// Close SQL
	err := d.SQL.Close()
	if err != nil {
		return errors.Wrap(err, "SQL DB Close")
	}

	d.Log.Debug(context.Background(), "SQL DB Closed")

	// Close NoSQL
	err = d.NoSQL.Disconnect(context.Background())
	if err != nil {
		return errors.Wrap(err, "NoSQL DB Close")
	}

	d.Log.Debug(context.Background(), "NoSQL DB Closed")

	// Close RMQ
	err = d.RMQ.Close(context.Background())
	if err != nil {
		return errors.Wrap(err, "RMQ DB Close")
	}

	d.Log.Debug(context.Background(), "RMQ DB Closed")

	return nil
}
