package main

import (
	"log"
	"github.com/vision-first/wegod/cmd/server/http/ginimpl"
)

func main() {
	if err := ginimpl.RunServer(); err != nil {
		log.Fatal(err)
	}
}
