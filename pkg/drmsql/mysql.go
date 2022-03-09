package drmsql

import (
	"context"
	"fmt"

	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/pkg/errors"
	"github.com/tsenart/nap"

	// Importing mysql driver dependencies
	_ "github.com/go-sql-driver/mysql"
)

// Config holds init dependencies for New DB
type Config struct {
	SQLConfig config.SQLConfig
	Log       drmlog.Logger
}

// GetDB open and pings connection to provided databases
// and return pointer DB struct with errors
func GetDB(conf Config) (*nap.DB, error) {
	dsns := getDataSourceName(conf.SQLConfig)

	db, err := nap.Open("mysql", dsns)
	if err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Open")
	}

	conf.Log.Info(context.Background(), "SQl Database Connected")

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "SQL Database Connection Testing")
	}

	conf.Log.Info(context.Background(), "SQl Database Connection Tested")

	return db, nil
}

func getDataSourceName(conf config.SQLConfig) (dsns string) {
	var (
		DBConf config.SQLConnConfig = conf[0]
		i      int
	)

	for i = 0; i < len(conf)-1; {
		dsns += fmt.Sprintf("%s:%s@tcp(%s:%d)/%s;", DBConf.User, DBConf.Password, DBConf.Host, DBConf.Port, DBConf.DB)
		DBConf = conf[i]
		i++
	}

	dsns += fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", DBConf.User, DBConf.Password, DBConf.Host, DBConf.Port, DBConf.DB)

	return
}
