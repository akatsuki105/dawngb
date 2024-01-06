package sched

type Sched struct {
	cycles  int64
	staging int64
	events  []*Event // Sorted by when (ascending)
}

func New() *Sched {
	return &Sched{}
}

func (s *Sched) Reset() {
	s.cycles = 0
	s.staging = 0
	s.events = make([]*Event, 0)
}

func (s *Sched) Cycle() int64 {
	return s.cycles + s.staging
}

func (s *Sched) Add(cycles int64) {
	s.staging += cycles
}

func (s *Sched) Schedule(event *Event, after int64) {
	event.when = after + s.Cycle()
	s.events = append(s.events, event)
	// Sort event by when
	for i := len(s.events) - 1; i > 0; i-- {
		if s.events[i].when < s.events[i-1].when {
			s.events[i], s.events[i-1] = s.events[i-1], s.events[i]
		} else {
			break
		}
	}
}

func (s *Sched) Cancel(event *Event) {
	for i, e := range s.events {
		if e == event {
			s.events = append(s.events[:i], s.events[i+1:]...)
			return
		}
	}
}

func (s *Sched) Commit() {
	s.cycles += s.staging
	s.staging = 0
	// Execute events
	for len(s.events) > 0 && s.events[0].when <= s.cycles {
		event := s.events[0]
		s.events = s.events[1:]
		event.Callback(s.cycles - event.when)
	}
}

func (s *Sched) UntilNextEvent() int64 {
	if len(s.events) == 0 {
		return 0
	}
	return s.events[0].when - s.cycles
}
