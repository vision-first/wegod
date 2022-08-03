package facades

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/log-go/impls/fmts"
	"github.com/995933447/log-go/impls/loggerwriters"
	simpletracectx "github.com/995933447/simpletrace/context"
	"github.com/995933447/std-go/print"
	"github.com/gin-gonic/gin"
	"github.com/vision-first/wegod/internal/app/ginimpl/enum"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"sync"
)

var (
	newLoggerMu sync.Mutex
	logger *log.Logger
)

func MustLogger() *log.Logger {
	if logger == nil {
		newLoggerMu.Lock()
		defer newLoggerMu.Unlock()
		if logger == nil {
			switch config.Conf.Log.Driver {
			case config.LogDriverFile:
				logger = log.NewLogger(MustNewFileLoggerWriter())
			default:
				logger = log.NewLogger(MustNewFileLoggerWriter())
			}
		}
	}

	return logger
}

func MustNewFileLoggerWriter() log.LoggerWriter {
	logger := loggerwriters.NewFileLoggerWriter(
		config.Conf.Log.Dir,
		config.Conf.Log.MaxFileSize,
		10,
		facades.CheckTimeToOpenNewFileHandlerForFileLogger(),
		100000,
	)

	go func() {
		if err := logger.Loop(); err != nil {
			panic(err)
		}
	}()

	logger.SetFormatter(NewGinSimpleTraceFormatter())

	return logger
}

type GinSimpleTraceFormatter struct {
	fmt *fmts.SimpleTraceFormatter
}

func NewGinSimpleTraceFormatter() *GinSimpleTraceFormatter {
	return &GinSimpleTraceFormatter{
		fmt: fmts.NewSimpleTraceFormatter(5, fmts.FormatText),
	}
}

func (f *GinSimpleTraceFormatter) Sprintf(ctx context.Context, level log.Level, stdoutColor print.Color, format string, args ...interface{}) (string, error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if traceCtx, ok := ginCtx.Get(enum.GinCtxKeySimpleTraceCtx); ok {
			return f.fmt.Sprintf(traceCtx.(*simpletracectx.Context), level, stdoutColor, format, args...)
		}
	}

	content, err := f.fmt.Sprintf(ctx, level, stdoutColor, format, args...)
	if err != nil {
		return "", err
	}

	return content, nil
}




