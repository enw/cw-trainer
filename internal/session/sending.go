package session

import "time"

// Sending stays in the same output state machine: reference and replay both
// cross the existing stop acknowledgment barrier before another output starts.
func (s *Session) sendingInput(in Input, now time.Time) ([]Effect, bool) {
	switch in.Kind {
	case BeginCapture, FinishCapture, ReplayAttempt, Reference, KeyEvent, CaptureClock, InputLost:
	default:
		return nil, false
	}
	if !s.SendBack {
		return nil, true
	}
	switch in.Kind {
	case BeginCapture:
		if s.Phase == Ready || s.Phase == Captured {
			s.Attempt.Begin(in.At)
			s.Phase = Capturing
		}
	case FinishCapture:
		if s.Phase == Capturing {
			_ = s.Attempt.Finish()
			s.Phase = Captured
		}
	case KeyEvent:
		if s.Phase == Capturing {
			s.Attempt.Feed(in.Event)
			if s.Attempt.Err != nil {
				s.Phase = Captured
			}
		}
	case CaptureClock:
		if s.Phase == Capturing {
			s.Attempt.Check(in.At)
			if s.Attempt.Err != nil {
				s.Phase = Captured
			}
		}
	case ReplayAttempt:
		if s.Phase == Captured || s.Phase == Ready {
			if _, err := s.Attempt.Timeline(); err == nil {
				return s.transition(Replay, now), true
			}
		}
	case Reference:
		if s.Attempt.Active {
			s.Attempt.Invalidate("Reference requested during capture; retry")
		}
		return s.transition(CW, now), true
	case InputLost:
		s.Attempt.Invalidate("MIDI connection lost; reconnect and start a fresh attempt")
		return s.transition(Captured, now), true
	}
	return nil, true
}
