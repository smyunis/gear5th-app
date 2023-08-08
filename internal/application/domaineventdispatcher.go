package application

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

var ApplicationDomainEventDispatcher DomainEventDispatcher

func init() {
	ApplicationDomainEventDispatcher = make(DomainEventDispatcher)
}

type DomainEventHandler func(any)

type DomainEventDispatcher map[string][]DomainEventHandler

func (d DomainEventDispatcher) AddHandler(eventname string, handler DomainEventHandler) {
	handlerEntries, ok := d[eventname]

	if !ok {
		d[eventname] = []DomainEventHandler{handler}
	} else {
		handlerEntries = append(handlerEntries, handler)
		d[eventname] = handlerEntries
	}
}

func (d DomainEventDispatcher) DispatchAsync(events ...shared.DomainEvents) {
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
