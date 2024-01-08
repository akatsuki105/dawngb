package sched

type Event struct {
	name     string // For debug
	Callback func(cyclesLate int64)
	when     int64
}

func NewEvent(name string, callback func(cyclesLate int64)) *Event {
	return &Event{
		name:     name,
		Callback: callback,
	}
}
