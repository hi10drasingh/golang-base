package drmsql

import (
	"context"
	"fmt"
	"time"

	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmtime"
	"github.com/pkg/errors"
	"github.com/tsenart/nap"

	// Importing mysql driver dependencies.
	_ "github.com/go-sql-driver/mysql"
)

// ConnConfig holds authentication details for a single sql server.
type ConnConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,number"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	DB       string `json:"db" validate:"required"`
}

// Config holds auth details of all servers.
type Config struct {
	Servers           []ConnConfig       `json:"servers" validate:"required"`
	ConnectionTimeout drmtime.CustomTime `json:"connectionTimeout" validate:"required"`
}

// GetDB open and pings connection to provided databases
// and return pointer DB struct with errors.
func GetDB(conf *Config, log drmlog.Logger) (*nap.DB, error) {
	dsns := getDataSourceName(conf.Servers, conf.ConnectionTimeout.Time)

	database, err := nap.Open("mysql", dsns)
	if err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Open")
	}

	if err := database.Ping(); err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Testing")
	}

	log.Debug(context.Background(), "SQl Database Connection Tested")

	return database, nil
}

func getDataSourceName(servers []ConnConfig, timeout time.Duration) (dsns string) {
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
