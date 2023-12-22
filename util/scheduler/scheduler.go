package scheduler

import (
	"fmt"
)

// Scheduler is used to synchronize between the components of the emulator.
//
// Scheduler's idea is from mGBA's one, so uses its timing mechanism for scheduling all kinds of stuff like how long it takes for a cartridge memory access or how long until the next time the audio engine needs to update.
type Scheduler struct {
	// Number of elapsed clocks of emulation core
	// Ja: エミュレーションコアの経過クロック数(すでにコミット(イベントを処理)したサイクル数)
	cycles int64

	// Relative offset from the start of the latest CPU step.
	//
	// Basically, This is pointer to cpu.Cycles
	// Ja: CPUステップ開始からの相対オフセット(まだコミットしていないサイクル数)
	relativeCycles *int64

	root *Event

	// How much time is left until the most recent event?
	//
	// *s.nextEvent = *s.relativeCycles + after(s.Schedule param)
	//
	// Basically, This is pointer to cpu.NextEvent
	// 次のイベントまでの残りさ時間(.cyclesからの相対時間)
	nextEvent *int64
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Reset(relativeCycles, nextEvent *int64) {
	s.cycles = 0
	s.root = nil
	s.SetCycles(relativeCycles, nextEvent)
}

// For serialization
func (s *Scheduler) SetCycles(relativeCycles, nextEvent *int64) {
	s.relativeCycles = relativeCycles
	s.nextEvent = nextEvent
}

// For serialization
func (s *Scheduler) MasterCycle() *int64 {
	return &s.cycles
}

// For serialization
func (s *Scheduler) RelativeCycle() int64 {
	return *s.relativeCycles
}

// Current cycles(with relativeCycles)
func (s *Scheduler) Cycle() int64 {
	return s.cycles + *s.relativeCycles
}

// This func is inspired by mgba's mTimingTick
func (s *Scheduler) Add(c int64) int64 {
	s.cycles += c
	masterCycles := s.cycles
	for s.root != nil {
		next := s.root
		nextWhen := next.when - masterCycles
		if nextWhen > 0 {
			return nextWhen
		}
		s.root = next.next
		next.Callback(-nextWhen)
	}
	return *s.nextEvent
}

// This func is inspired by mgba's mTimingSchedule
func (s *Scheduler) Schedule(event *Event, after int64) {
	after += *s.relativeCycles
	event.when = after + s.cycles
	if after < *s.nextEvent {
		*s.nextEvent = after
	}

	previous := &s.root
	next := s.root
	priority := event.priority
	for next != nil {
		nextWhen := next.when - s.cycles
		if nextWhen > after || (nextWhen == after && next.priority > priority) {
			break
		}

		previous = &next.next
		next = next.next
	}

	event.next = next
	*previous = event
}

func (s *Scheduler) ReSchedule(e *Event, after int64) {
	s.Deschedule(e)
	s.Schedule(e, after)
}

func (s *Scheduler) ScheduleAbs(e *Event, when int64) {
	s.Schedule(e, when-s.Cycle())
}

func (s *Scheduler) Deschedule(event *Event) {
	previous := &s.root
	next := s.root
	for next != nil {
		if next == event {
			*previous = next.next
			return
		}
		previous = &next.next
		next = next.next
	}
}

func (s *Scheduler) Until(e *Event) int64 {
	return e.when - s.cycles - *s.relativeCycles
}

// This func is inspired by mgba's mTimingIsScheduled
func (s *Scheduler) Scheduled(e *Event) bool {
	next := s.root
	if s.root == nil {
		return false
	}

	for next != nil {
		if next == e {
			return true
		}
		next = next.next
	}
	return false
}

func (s *Scheduler) String() string {
	result := ""
	event := s.root
	for event != nil {
		result += fmt.Sprintf("%s:%d->", event.name, event.when)
		event = event.next
	}
	return result
}

func (s *Scheduler) Events() []*Event {
	result := []*Event{}

	next := s.root
	for next != nil {
		result = append(result, next)
		next = next.next
	}
	return result
}
