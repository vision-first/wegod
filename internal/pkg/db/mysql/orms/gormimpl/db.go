package gormimpl

import (
	"fmt"
	"github.com/995933447/log-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type InitDBConfig struct {
	UserName string
	Password string
	Host string
	Port uint32
	Database string
	Charset string
	MaxConns int
	MaxIdleConns int
	LoggerConfig *DBLoggerConfig
}

type DBLoggerConfig struct {
	LogLevel logger.LogLevel
	BaseLogger *log.Logger
	SlowThreshold time.Duration
}

func InitDB(conf *InitDBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Charset,
		)
	gormConf := &gorm.Config{}
	if conf.LoggerConfig != nil {
		gormConf.Logger = NewLogger(conf.LoggerConfig.LogLevel, conf.LoggerConfig.BaseLogger, conf.LoggerConfig.SlowThreshold)
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConf)
	if err != nil {
		return nil, err
	}
	if conf.MaxConns > 0 || conf.MaxIdleConns > 0 {
		if sqlDB, err := db.DB(); err != nil {
			return nil, err
		} else {
			if conf.MaxConns > 0 {
				sqlDB.SetMaxIdleConns(conf.MaxConns)
			}
			if conf.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
			}
		}
	}

	return db, nil
}