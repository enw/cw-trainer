// Package session owns playback policy. Only output effects cross into workers.
package session

import "time"

type Phase string

const (
	Idle     Phase = "Stopped"
	CW       Phase = "CW"
	Speech   Phase = "Speaking"
	Gap      Phase = "Gap"
	Paused   Phase = "Paused"
	Stopping Phase = "Stopping"
	Failed   Phase = "Error"
)

type Kind int

const (
	Start Kind = iota
	Next
	Previous
	TogglePause
	Speak
	Back
	Quit
	Done
	Stopped
	Error
)

type Input struct {
	Kind Kind
	ID   uint64
	Err  error
}
type EffectKind int

const (
	PlayCW EffectKind = iota
	PlaySpeech
	Wait
	StopOutput
	Exit
)

type Effect struct {
	Kind     EffectKind
	ID       uint64
	Duration time.Duration
}

type Session struct {
	Phase               Phase
	Line                int
	Count               int
	Err                 error
	lines               int
	id                  uint64
	pending             Phase
	resume              Phase
	gapRemaining        time.Duration
	gapUntil            time.Time
	speechReturnsPaused bool
	quitting            bool
}

func New(lines int) *Session { return &Session{Phase: Idle, lines: lines, resume: CW} }

// Update is serialized by the UI event loop; time is supplied for deterministic
// gap tests. Output completion is accepted only for the active generation.
func (s *Session) Update(in Input, now time.Time) []Effect {
	if s.quitting && in.Kind != Stopped && in.Kind != Done && in.Kind != Error {
		return nil
	}
	switch in.Kind {
	case Done, Stopped, Error:
		return s.outputEvent(in, now)
	case Quit:
		s.quitting = true
		return s.transition(Idle, now)
	case Back:
		s.speechReturnsPaused = false
		return s.transition(Idle, now)
	case Start:
		if s.lines == 0 {
			return nil
		}
		s.Err = nil
		s.Count = 0
		s.resume = CW
		return s.transition(CW, now)
	case Next, Previous:
		step := 1
		if in.Kind == Previous {
			step = -1
		}
		target := min(max(s.Line+step, 0), s.lines-1)
		if s.lines == 0 || target == s.Line {
			return nil
		}
		s.Line = target
		s.Count = 0
		s.resume = CW
		paused := s.Phase == Paused || (s.Phase == Stopping && s.pending == Paused) || s.speechReturnsPaused
		s.speechReturnsPaused = false
		if paused {
			return s.transition(Paused, now)
		}
		if s.Phase == Idle || s.Phase == Failed {
			return nil
		}
		return s.transition(CW, now)
	case TogglePause:
		return s.pause(now)
	case Speak:
		if s.lines == 0 || s.Phase == Idle || s.Phase == Failed {
			return nil
		}
		s.speechReturnsPaused = s.Phase == Paused || (s.Phase == Stopping && s.pending == Paused) || s.speechReturnsPaused
		s.Count = 0
		s.resume = CW
		return s.transition(Speech, now)
	}
	return nil
}

func (s *Session) transition(target Phase, now time.Time) []Effect {
	s.pending = target
	if s.Phase == Stopping {
		return nil
	}
	if s.Phase == CW || s.Phase == Speech || s.Phase == Gap {
		s.Phase = Stopping
		return []Effect{{Kind: StopOutput, ID: s.id}}
	}
	return s.begin(target, now)
}
func (s *Session) begin(target Phase, now time.Time) []Effect {
	s.Phase = target
	if target == Idle && s.quitting {
		return []Effect{{Kind: Exit}}
	}
	var kind EffectKind
	switch target {
	case CW:
		kind = PlayCW
	case Speech:
		kind = PlaySpeech
	case Gap:
		kind = Wait
		s.gapUntil = now.Add(s.gapRemaining)
	default:
		return nil
	}
	s.id++
	return []Effect{{Kind: kind, ID: s.id, Duration: s.gapRemaining}}
}
func (s *Session) outputEvent(in Input, now time.Time) []Effect {
	if in.ID != s.id {
		return nil
	}
	if s.Phase == Stopping {
		if in.Kind == Stopped {
			return s.begin(s.pending, now)
		}
		return nil
	}
	if in.Kind == Stopped || !(s.Phase == CW || s.Phase == Speech || s.Phase == Gap) {
		return nil
	}
	if in.Kind == Error {
		s.Phase = Failed
		s.Err = in.Err
		return nil
	}
	switch s.Phase {
	case CW:
		s.Count++
	case Speech:
		if s.speechReturnsPaused {
			s.speechReturnsPaused = false
			s.resume = CW
			return s.begin(Paused, now)
		}
	case Gap:
		return s.begin(CW, now)
	}
	s.gapRemaining = 2 * time.Second
	return s.begin(Gap, now)
}
func (s *Session) pause(now time.Time) []Effect {
	if s.Phase == Idle || s.Phase == Failed {
		return nil
	}
	if s.Phase == Paused {
		return s.begin(s.resume, now)
	}
	if s.Phase == Stopping {
		if s.pending == Paused {
			s.pending = s.resume
		} else {
			s.resume = s.pending
			s.pending = Paused
		}
		return nil
	}
	s.resume = s.Phase
	if s.Phase == Gap {
		s.gapRemaining = max(0, s.gapUntil.Sub(now))
	}
	return s.transition(Paused, now)
}
