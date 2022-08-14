package main

import (
	"fmt"
	"github.com/995933447/log-go"
	"github.com/995933447/log-go/impls/loggerwriters"
	"github.com/995933447/std-go/scan"
	"github.com/vision-first/wegod/internal/pkg/config"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"os"
	"path"
	"strings"
	"sync"
)

type (
	consoleHandleFunc func()
	consoleHandlerMetadata struct {
		Name  string
		Usage string
		Handler consoleHandleFunc
	}
)

var consoleHandlerMetadataList []*consoleHandlerMetadata
var mu sync.Mutex

func Register(name, usage string, consoleHandleFunc consoleHandleFunc) {
	mu.Lock()
	defer mu.Unlock()
	consoleHandlerMetadataList = append(consoleHandlerMetadataList, &consoleHandlerMetadata{Name: name, Usage: usage, Handler: consoleHandleFunc})
}

func showUsage() {
	fmt.Println()
	for _, medata := range consoleHandlerMetadataList {
		fmt.Printf("\t-f %s %s\n", medata.Name, medata.Usage)
	}
	fmt.Println()
	os.Exit(1)
}

func Run() {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-h" || os.Args[i] == "--help" {
			showUsage()
		}
	}

	callHandlerName := scan.OptStr("f")
	if callHandlerName == "" {
		fmt.Printf("\nmissed -f option\n\n")
		showUsage()
	}

	name := path.Base(strings.Replace(os.Args[0], "\\", "/", -1))
	if strings.HasPrefix(name, "./") {
		name = name[2:]
	}

	var handler consoleHandleFunc
	for _, metadata := range consoleHandlerMetadataList {
		if strings.ToLower(metadata.Name) == strings.ToLower(callHandlerName) {
			handler = metadata.Handler
			break
		}
	}

	if handler == nil {
		fmt.Printf("\nnot found handler %s\n\n", callHandlerName)
		showUsage()
	}

	handler()
}

func init() {
	loggerWriter := loggerwriters.NewFileLoggerWriter(
		config.Conf.Log.Dir,
		config.Conf.Log.MaxFileSize,
		10,
		facades.CheckTimeToOpenNewFileHandlerForFileLogger(),
		100000,
	)
	go func() {
		if err := loggerWriter.Loop(); err != nil {
			panic(err)
		}
	}()
	logger = log.NewLogger(loggerWriter)
}

func main()  {
	Register("GenApi", "-a Api -m Method", GenApi)
	Register("GenService", "-s Service", GenService)
	Run()
}
