package main

import (
	"fmt"
	"github.com/vision-first/wegod/tools/autogen/apigen"
	"path/filepath"
)

func main()  {
	fmt.Println(filepath.Abs("./internal/pkg/api"))
	err := apigen.GenApi(
		"buddha",
		"PleaseBuddha",
		"./internal/pkg/api",
		"./tools/autogen/templates")
	if err != nil {
		panic(err)
	}
}
