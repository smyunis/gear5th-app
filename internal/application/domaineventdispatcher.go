package application

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"


var x EventDispatcher = appEventDispatcher

type EventDispatcher interface {
	AddHandler(eventname string, handler EventHandler)
	DispatchAsync(events ...shared.Events)
}

var appEventDispatcher InMemoryEventDispatcher

func init() {
	appEventDispatcher = make(InMemoryEventDispatcher)
}

func NewAppEventDispatcher() InMemoryEventDispatcher {
	return appEventDispatcher
}

type EventHandler func(any)
type InMemoryEventDispatcher map[string][]EventHandler

func (d InMemoryEventDispatcher) AddHandler(eventname string, handler EventHandler) {
	handlerEntries, ok := d[eventname]

	if !ok {
		d[eventname] = []EventHandler{handler}
	} else {
		handlerEntries = append(handlerEntries, handler)
		d[eventname] = handlerEntries
	}
}

func (d InMemoryEventDispatcher) DispatchAsync(events ...shared.Events) {
	for _, eventSet := range events {
		for eventName, eventPayloads := range eventSet {
			handlers, ok := d[eventName]
			if !ok {
				continue
			}
			for _, handler := range handlers {
				for _, event := range eventPayloads {
					go handler(event)
				}
			}
		}
	}
}

