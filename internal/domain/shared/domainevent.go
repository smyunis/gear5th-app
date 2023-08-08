package shared

type DomainEvents map[string][]any

func (d *DomainEvents) Emit(eventname string, event any) {
	eventEntries, ok := (*d)[eventname]
	if !ok {
		(*d)[eventname] = []any{event}
	} else {
		eventEntries = append(eventEntries, event)
		(*d)[eventname] = eventEntries
	}
}
