// Package keying captures measured marks and spaces without decoding them.
package keying

import (
	"cw-looper/internal/morse"
	"fmt"
	"time"
)

const MaxDuration = 2 * time.Minute
const MaxMarks = 1024

type Event struct {
	At   time.Duration
	Note byte
	Down bool
}
type Mark struct{ Start, End time.Duration }
type Attempt struct {
	Marks               []Mark
	Err                 error
	Active              bool
	armed, last, downAt time.Duration
	down                bool
	note                byte
	complete            bool
}

func (a *Attempt) Begin(at time.Duration) { *a = Attempt{Active: true, armed: at, last: at} }
func (a *Attempt) Invalidate(reason string) {
	a.Err = fmt.Errorf("%s", reason)
	a.Active = false
	a.complete = false
	a.down = false
}
func (a *Attempt) Check(now time.Duration) {
	if a.Active && now-a.armed > MaxDuration {
		a.Invalidate("Attempt exceeded two minutes; retry")
	}
}
func (a *Attempt) Feed(event Event) {
	if !a.Active || event.At < a.armed {
		return
	}
	a.Check(event.At)
	if !a.Active {
		return
	}
	if event.At < a.last {
		a.Invalidate("Out-of-order keying timestamps; retry")
		return
	}
	a.last = event.At
	if event.Down {
		if a.down {
			a.Invalidate("Overlapping key-down events; release pads and retry")
			return
		}
		a.down = true
		a.downAt = event.At
		a.note = event.Note
		return
	}
	if !a.down || event.Note != a.note || event.At <= a.downAt {
		a.Invalidate("Incomplete keying event; release pads and retry")
		return
	}
	a.Marks = append(a.Marks, Mark{a.downAt, event.At})
	a.down = false
	if len(a.Marks) > MaxMarks {
		a.Invalidate("Attempt exceeded 1024 marks; retry")
	}
}
func (a *Attempt) Finish() error {
	if a.Err != nil {
		return a.Err
	}
	if !a.Active {
		return fmt.Errorf("No active attempt")
	}
	a.Active = false
	if a.down {
		a.Invalidate("Attempt ended while key was down; release pads and retry")
	} else if len(a.Marks) == 0 {
		a.Invalidate("No keying received; retry")
	}
	a.complete = a.Err == nil
	return a.Err
}
func (a *Attempt) Timeline() (morse.Timeline, error) {
	var timeline morse.Timeline
	if !a.complete || a.Err != nil || len(a.Marks) == 0 {
		return timeline, fmt.Errorf("No complete attempt to replay")
	}
	origin := a.Marks[0].Start
	sample := func(d time.Duration) int64 {
		return (int64(d)*morse.SampleRate + int64(time.Second)/2) / int64(time.Second)
	}
	for _, mark := range a.Marks {
		start, end := sample(mark.Start-origin), sample(mark.End-origin)
		if start > timeline.Samples {
			timeline.Segments = append(timeline.Segments, morse.Segment{Start: timeline.Samples, End: start})
		}
		if end <= start {
			return morse.Timeline{}, fmt.Errorf("Captured mark is shorter than one audio sample")
		}
		timeline.Segments = append(timeline.Segments, morse.Segment{Start: start, End: end, Tone: true})
		timeline.Samples = end
	}
	return timeline, nil
}
