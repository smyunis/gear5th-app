package testdoubles

import (
	"gitlab.com/gear5th/gear5th-app/internal/application"
	"gitlab.com/gear5th/gear5th-app/internal/domain/shared"
)

type LocalizedEventDispatcher struct {
	dispatcher map[string][]application.EventHandler
}

func (d *LocalizedEventDispatcher) AddHandler(eventname string, handler application.EventHandler) {
	handlerEntries, ok := d.dispatcher[eventname]

	if !ok {
		d.dispatcher[eventname] = []application.EventHandler{handler}
	} else {
		handlerEntries = append(handlerEntries, handler)
		d.dispatcher[eventname] = handlerEntries
	}
}

func (d *LocalizedEventDispatcher) DispatchAsync(events ...shared.Events) {
	for _, eventSet := range events {
		for eventName, eventPayloads := range eventSet {
			handlers, ok := d.dispatcher[eventName]
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
