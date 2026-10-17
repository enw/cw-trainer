package session

import (
	"cw-looper/internal/library"
	"time"
)

func (s *Session) resetUnit() {
	s.portion = 0
	s.unitInterval = 0
	s.unitComplete = false
	s.onDemand = false
	s.afterGap = CW
}

func (s *Session) completePortion() {
	s.afterGap = CW
	if s.onDemand {
		s.resetUnit()
		return
	}
	s.portion++
	switch s.Settings.Style {
	case "sandwich":
		if s.portion == 1 {
			s.afterGap = Speech
		}
		s.unitComplete = s.portion == 3
	case "occasional":
		if s.portion == s.unitInterval {
			s.afterGap = Speech
		}
		s.unitComplete = s.portion > s.unitInterval
		if s.unitComplete {
			s.Count = 0
		}
	default:
		s.unitComplete = true
	}
}

func (s *Session) finishUnit() {
	s.advanceUnit()
	s.resetUnit()
}

// Configure applies display immediately and numeric audio settings at the next
// CW. Policy changes cancel the current unit through the same stop barrier.
func (s *Session) Configure(settings library.Settings, now time.Time) []Effect {
	old := s.Settings
	s.Settings = settings
	if old.Display != settings.Display {
		s.Revealed = false
	}
	if old.Style != settings.Style || old.Sequence != settings.Sequence || old.LineAdvancement != settings.LineAdvancement {
		return s.navigate(now)
	}
	return nil
}
