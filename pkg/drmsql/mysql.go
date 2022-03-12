package drmsql

import (
	"context"
	"fmt"
	"time"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
	"github.com/tsenart/nap"

	// Importing mysql driver dependencies.
	_ "github.com/go-sql-driver/mysql"
)

// Config holds init dependencies for New DB.
type Config struct {
	SQLConfig config.SQLConfig
	Log       drmlog.Logger
}

// GetDB open and pings connection to provided databases
// and return pointer DB struct with errors.
func GetDB(conf *Config) (*nap.DB, error) {
	dsns := getDataSourceName(conf.SQLConfig.Servers, conf.SQLConfig.ConnectionTimeout.Time)

	database, err := nap.Open("mysql", dsns)
	if err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Open")
	}

	if err := database.Ping(); err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Testing")
	}

	conf.Log.Info(context.Background(), "SQl Database Connection Tested")

	return database, nil
}

func getDataSourceName(servers []config.SQLConnConfig, timeout time.Duration) (dsns string) {
	var (
		DBConf = servers[0]
		index  int
		format = "%s:%s@tcp(%s:%d)/%s?timeout=%s"
	)

	for index = 0; index < len(servers)-1; {
		dsns += fmt.Sprintf(format+";", DBConf.User, DBConf.Password, DBConf.Host, DBConf.Port, DBConf.DB, timeout.String())

		DBConf = servers[index]

		index++
	}

	dsns += fmt.Sprintf(format, DBConf.User, DBConf.Password, DBConf.Host, DBConf.Port, DBConf.DB, timeout.String())

	return dsns
}
