package main

import (
	"context"
	simpletracectx "github.com/995933447/simpletrace/context"
)

func NewCtx() *simpletracectx.Context {
	return simpletracectx.New("wegod-console", context.Background(), "", "")
}