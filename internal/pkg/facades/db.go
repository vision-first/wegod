package facades

import (
	"github.com/995933447/log-go"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
	"wegod/internal/pkg/config"
	"wegod/internal/pkg/db/mysql/orms/gormimpl"
)

var (
	newGormDBMu sync.Mutex
	gormDB *gorm.DB
)

func MustGormDB(baseLogger *log.Logger) *gorm.DB {
	if gormDB == nil {
		newGormDBMu.Lock()
		defer newGormDBMu.Unlock()
		var err error
		gormDB, err = gormimpl.InitDB(&gormimpl.InitDBConfig{
			UserName: config.Conf.DB.Mysql.UserName,
			Password: config.Conf.DB.Mysql.Password,
			Database: config.Conf.DB.Mysql.Database,
			Host: config.Conf.DB.Mysql.Host,
			Port: config.Conf.DB.Mysql.Port,
			Charset: config.Conf.Mysql.Charset,
			MaxIdleConns: config.Conf.Mysql.MaxIdleConns,
			MaxConns: config.Conf.Mysql.MaxConns,
			LoggerConfig: &gormimpl.DBLoggerConfig{
				LogLevel: transConfLogLevelToGorm(config.Conf.DB.Mysql.Log.Level),
				SlowThreshold: time.Second * time.Duration(config.Conf.Mysql.Log.SlowThreshold),
				BaseLogger: baseLogger,
			},
		})
		if err != nil {
			panic(err)
		}
	}
	return gormDB
}

func transConfLogLevelToGorm(confLogLevel string) logger.LogLevel {
	switch confLogLevel {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}