package main

import (
	"github.com/995933447/log-go"
	"github.com/995933447/log-go/impls/loggerwriters"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

var logger = log.NewLogger(loggerwriters.NewFileLoggerWriter(
	config.Conf.Log.Dir,
	config.Conf.Log.MaxFileSize,
	10,
	facades.CheckTimeToOpenNewFileHandlerForFileLogger(),
	100000,
	))
