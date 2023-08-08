package shared

type Events map[string][]any

func (d *Events) Emit(eventname string, event any) {
	eventEntries, ok := (*d)[eventname]
	if !ok {
		(*d)[eventname] = []any{event}
	} else {
		eventEntries = append(eventEntries, event)
		(*d)[eventname] = eventEntries
	}
}
