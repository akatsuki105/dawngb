package scheduler

// For debug, string is used
type EventName string

type Event struct {
	name     EventName
	Callback func(cyclesLate int64)
	when     int64
	priority uint
	next     *Event
}

func NewEvent(name EventName, callback func(cyclesLate int64), prio uint) *Event {
	return &Event{
		name:     name,
		Callback: callback,
		priority: prio,
	}
}

// Name returns event name
func (e *Event) Name() EventName {
	return e.name
}

// When this event is triggerd
func (e *Event) When() *int64 {
	return &e.when
}

// Event's priority
func (e *Event) Priority() uint {
	return e.priority
}
