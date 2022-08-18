package facades

import (
	"github.com/995933447/eventobserver"
	"github.com/995933447/log-go"
	"sync"
)

var (
	eventDispatcher *eventobserver.Dispatcher
	newEventDispatcherMu sync.Mutex
)

func EventDispatcher(logger *log.Logger) *eventobserver.Dispatcher {
	if eventDispatcher != nil {
		return eventDispatcher
	}

	newEventDispatcherMu.Lock()
	defer newEventDispatcherMu.Unlock()

	if eventDispatcher != nil {
		return eventDispatcher
	}

	eventDispatcher = eventobserver.NewDispatcher(logger)

	return eventDispatcher
}