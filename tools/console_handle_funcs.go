package main

import (
	"github.com/995933447/std-go/scan"
	"github.com/vision-first/wegod/tools/autogen/apigen"
	"github.com/vision-first/wegod/tools/autogen/servicegen"
)

func GenApi() {
	api := scan.OptStr("a")
	method := scan.OptStr("m")
	err := apigen.GenApi(
		api,
		method,
		"github.com/vision-first/wegod/internal/pkg/api",
		"../internal/pkg/api",
	)
	if err != nil {
		logger.Error(NewCtx(), err)
		panic(err)
	}
}

func GenService() {
	srv := scan.OptStr("s")
	err := servicegen.GenService(
		srv,
		"../internal/pkg/services",
	)
	if err != nil {
		logger.Error(NewCtx(), err)
		panic(err)
	}
}

func GenServicePageQuery() {
	srv := scan.OptStr("s")
	method := scan.OptStr("m")
	dataModel := scan.OptStr("mod")
	err := servicegen.GenServiceOptionStreamPageQuery(
		srv,
		method,
		dataModel,
		"github.com/vision-first/wegod/internal/pkg/datamodel/models",
		"github.com/vision-first/wegod/internal/pkg/facades",
		"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl",
		"../internal/pkg/services",
	)
	if err != nil {
		logger.Error(NewCtx(), err)
		panic(err)
	}
}
