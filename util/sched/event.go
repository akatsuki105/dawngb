package sched

type Event struct {
	name     string // For debug
	Callback func(cyclesLate int64)
	when     int64
	priority uint
}

func NewEvent(name string, callback func(cyclesLate int64), prio uint) *Event {
	return &Event{
		name:     name,
		Callback: callback,
		priority: prio,
	}
}
