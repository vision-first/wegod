package eventobserver

import (
	"context"
	"github.com/995933447/log-go"
	"sync"
)

type Event struct {
	name string
	data interface{}
}

func NewEvent(name string, data interface{}) *Event {
	return &Event{
		name: name,
		data: data,
	}
}

func (e *Event) GetName() string {
	return e.name
}

func (e *Event) GetData() interface{} {
	return e.data
}

type HandleEventFunc func(ctx context.Context, event Event) error

type Listener struct {
	name string
	eventHandler HandleEventFunc
}

func NewListener(name string, eventHandler HandleEventFunc) *Listener {
	return &Listener{
		name: name,
		eventHandler: eventHandler,
	}
}

type Dispatcher struct {
	eventNameToListenerNamesMap map[string][]string
	listenerMap map[string]*Listener
	mu sync.RWMutex
	logger *log.Logger
}

func NewDispatcher(logger *log.Logger) *Dispatcher {
	return &Dispatcher{
		eventNameToListenerNamesMap: make(map[string][]string),
		listenerMap: make(map[string]*Listener),
		logger: logger,
	}
}

func (d *Dispatcher) AddListener(eventName string, listener *Listener)  {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.eventNameToListenerNamesMap[eventName] = append(d.eventNameToListenerNamesMap[eventName], listener.name)
	d.listenerMap[listener.name] = listener
}

func (d *Dispatcher) DelListener(eventName string, listener *Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	listenerNames := d.eventNameToListenerNamesMap[eventName]
	for i, name := range listenerNames {
		if name != listener.name {
			continue
		}

		listenerNamesSize := len(listenerNames)
		newListenerNames := make([]string, listenerNamesSize - 1)
		copy(newListenerNames, listenerNames[0:i])
		if i + 1 < listenerNamesSize {
			copy(newListenerNames, listenerNames[i + 1:])
		}
		listenerNames = newListenerNames
		break
	}
	d.eventNameToListenerNamesMap[eventName] = listenerNames
	delete(d.listenerMap, listener.name)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event Event) {
	listenerNames := d.eventNameToListenerNamesMap[event.name]
	if len(listenerNames) == 0 {
		return
	}

	for _, name := range listenerNames {
		listener, ok := d.listenerMap[name]
		if !ok {
			d.logger.Warnf(ctx, "listener:%s not found.", name)
			continue
		}

		if err := listener.eventHandler(ctx, event); err != nil {
			d.logger.Errorf(ctx, "listener:%s handle event occur error, stop handle listener chain.", listener.name)
			break
		}
	}
}