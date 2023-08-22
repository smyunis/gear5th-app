package application

import "gitlab.com/gear5th/gear5th-app/internal/domain/shared"

var ApplicationEventDispatcher EventDispatcher

func init() {
	ApplicationEventDispatcher = make(EventDispatcher)
}

type EventHandler func(any)

type EventDispatcher map[string][]EventHandler

func (d EventDispatcher) AddHandler(eventname string, handler EventHandler) {
	handlerEntries, ok := d[eventname]

	if !ok {
		d[eventname] = []EventHandler{handler}
	} else {
		handlerEntries = append(handlerEntries, handler)
		d[eventname] = handlerEntries
	}
}

func (d EventDispatcher) DispatchAsync(events ...shared.Events) {
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
