package facades

import (
	"github.com/995933447/log-go"
	"github.com/995933447/eventobserver"
	"sync"
)

var (
	eventDispatcher *eventobserver.Dispatcher
	newEventDispatcherMu sync.Mutex
)

func EventDispatcher(logger *log.Logger) *eventobserver.Dispatcher {
	if eventDispatcher == nil {
		newEventDispatcherMu.Lock()
		defer newEventDispatcherMu.Unlock()
		if eventDispatcher == nil {
			eventDispatcher = eventobserver.NewDispatcher(logger)
		}
	}
	return eventDispatcher
}