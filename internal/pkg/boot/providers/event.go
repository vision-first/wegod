package providers

import (
	"github.com/995933447/eventobserver"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/boot"
	"github.com/vision-first/wegod/internal/pkg/facades"
)

type EventProvider struct {
	eventNameToListenersMap map[string][]*eventobserver.Listener
	logger *log.Logger
}

var _ boot.ServiceProvider = (*EventProvider)(nil)

func NewEventProvider(eventNameToListenersMap map[string][]*eventobserver.Listener, logger *log.Logger) *EventProvider {
	return &EventProvider{
		eventNameToListenersMap: eventNameToListenersMap,
		logger: logger,
	}
}

func (e *EventProvider) Boot() error {
	dispatcher := facades.EventDispatcher(e.logger)
	for eventName, listeners := range e.eventNameToListenersMap {
		for _, listener := range listeners {
			dispatcher.AddListener(eventName, listener)
		}
	}
	return nil
}
