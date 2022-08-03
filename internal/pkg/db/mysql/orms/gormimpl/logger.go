package gormimpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/995933447/log-go"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Logger struct {
	logLevel logger.LogLevel
	baseLogger *log.Logger
	slowThreshold time.Duration
}

func NewLogger(level logger.LogLevel, baseLogger *log.Logger, slowThreshold time.Duration) *Logger {
	return &Logger{
		logLevel: level,
		baseLogger: baseLogger,
		slowThreshold: slowThreshold,
	}
}

const (
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface  {
	l.logLevel = level
	return l
}

func (l *Logger) Info(ctx context.Context, format string,  args ...interface{}) {
	if l.logLevel >= logger.Info {
		l.baseLogger.Infof(ctx, format, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, format string,  args ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.baseLogger.Warnf(ctx, format, args...)
	}
}

func (l *Logger) Error(ctx context.Context, format string,  args ...interface{}) {
	if l.logLevel >= logger.Error {
		l.baseLogger.Errorf(ctx, format, args...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		if rows == -1 {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.logLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.baseLogger.Infof(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.baseLogger.Infof(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
