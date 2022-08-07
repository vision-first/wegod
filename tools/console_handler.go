package main

import (
	"github.com/995933447/std-go/scan"
	"github.com/vision-first/wegod/tools/autogen/apigen"
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
		panic(err)
	}
}
