// Package session owns playback policy. Only output effects cross into workers.
package session

import (
	"cw-looper/internal/library"
	"time"
)

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
	NextEntry
	PreviousEntry
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
	Settings library.Settings
	Kind     EffectKind
	ID       uint64
	Duration time.Duration
}

type Session struct {
	Settings, ActiveSettings library.Settings
	Phase                    Phase
	Line                     int
	Entry                    int
	entryLengths             []int
	Count                    int
	Err                      error
	lines                    int
	id                       uint64
	pending                  Phase
	resume                   Phase
	gapRemaining             time.Duration
	gapUntil                 time.Time
	speechReturnsPaused      bool
	quitting                 bool
}

func New(lines int) *Session {
	return &Session{Settings: library.Defaults(), ActiveSettings: library.Defaults(), Phase: Idle, lines: lines, entryLengths: []int{lines}, resume: CW}
}

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
		return s.navigate(now)
	case NextEntry, PreviousEntry:
		step := 1
		if in.Kind == PreviousEntry {
			step = -1
		}
		target := min(max(s.Entry+step, 0), len(s.entryLengths)-1)
		if len(s.entryLengths) == 0 || target == s.Entry {
			return nil
		}
		s.Entry = target
		s.lines = s.entryLengths[target]
		s.Line = 0
		return s.navigate(now)
	case TogglePause:
		return s.pause(now)
	case Speak:
		if s.lines == 0 || s.Phase == Idle || s.Phase == Failed {
			return nil
		}
		s.speechReturnsPaused = s.wantsPause()
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
		s.ActiveSettings = s.Settings
	case Speech:
		kind = PlaySpeech
	case Gap:
		kind = Wait
		s.gapUntil = now.Add(s.gapRemaining)
	default:
		return nil
	}
	s.id++
	return []Effect{{Kind: kind, ID: s.id, Duration: s.gapRemaining, Settings: s.ActiveSettings}}
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
	s.gapRemaining = time.Duration(s.Settings.PauseSeconds * float64(time.Second))
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

// Select is the single-exchange form of SelectEntries.
func (s *Session) Select(lines, index int) { s.SelectEntries([]int{lines}, 0, index) }

// SelectEntries configures a stopped session without reusing generation IDs.
// The session owns both line and top-level navigation, including pending intent.
func (s *Session) SelectEntries(lengths []int, entry, index int) {
	if s.Phase != Idle && s.Phase != Failed {
		return
	}
	s.entryLengths = append([]int(nil), lengths...)
	s.Entry = min(max(entry, 0), max(0, len(lengths)-1))
	s.lines = 0
	if len(lengths) > 0 {
		s.lines = lengths[s.Entry]
	}
	s.Line = min(max(index, 0), max(0, s.lines-1))
	s.Count = 0
	s.Err = nil
	s.Phase = Idle
	s.resume = CW
	s.speechReturnsPaused = false
}
func (s *Session) wantsPause() bool {
	return s.Phase == Paused || (s.Phase == Stopping && s.pending == Paused) || s.speechReturnsPaused
}
func (s *Session) navigate(now time.Time) []Effect {
	paused := s.wantsPause()
	s.Count = 0
	s.resume = CW
	s.speechReturnsPaused = false
	if paused {
		return s.transition(Paused, now)
	}
	if s.Phase == Idle || s.Phase == Failed {
		return nil
	}
	return s.transition(CW, now)
}
