package session

import (
	"cw-looper/internal/keying"
	"testing"
	"time"
)

func TestSendBackCapturesThenReplaysAndIgnoresAutomaticProgression(t *testing.T) {
	s := New(2)
	s.SendBack = true
	s.Settings.Sequence = "cycle"
	s.Settings.Style = "sandwich"
	s.Settings.LineAdvancement = "automatic"
	now := time.Now()
	fx := s.Update(Input{Kind: Start}, now)[0]
	if got := s.Update(Input{Kind: Done, ID: fx.ID}, now); len(got) != 0 || s.Phase != Ready {
		t.Fatal("reference did not stop ready")
	}
	s.Update(Input{Kind: BeginCapture, At: time.Second}, now)
	for _, e := range []keying.Event{{At: time.Second, Note: 1, Down: true}, {At: 1080 * time.Millisecond, Note: 1}} {
		s.Update(Input{Kind: KeyEvent, Event: e}, now)
	}
	s.Update(Input{Kind: FinishCapture}, now)
	fx = s.Update(Input{Kind: ReplayAttempt}, now)[0]
	if fx.Kind != PlayAttempt || fx.Timeline.Samples != 3840 {
		t.Fatal("did not emit measured replay")
	}
	s.Update(Input{Kind: Done, ID: fx.ID}, now)
	if s.Phase != Captured || s.Line != 0 {
		t.Fatal("automatic practice skipped learner turn")
	}
}

func TestSendingRejectsEarlyInputAndStopsReplayBeforeNavigation(t *testing.T) {
	s := New(2)
	s.SendBack = true
	now := time.Now()
	fx := s.Update(Input{Kind: Start}, now)[0]
	s.Update(Input{Kind: KeyEvent, Event: keying.Event{At: time.Second, Down: true}}, now)
	s.Update(Input{Kind: Done, ID: fx.ID}, now)
	s.Update(Input{Kind: BeginCapture, At: 2 * time.Second}, now)
	s.Update(Input{Kind: KeyEvent, Event: keying.Event{At: time.Second, Down: true}}, now)
	s.Update(Input{Kind: KeyEvent, Event: keying.Event{At: 2 * time.Second, Down: true}}, now)
	s.Update(Input{Kind: KeyEvent, Event: keying.Event{At: 2080 * time.Millisecond}}, now)
	s.Update(Input{Kind: FinishCapture}, now)
	fx = s.Update(Input{Kind: ReplayAttempt}, now)[0]
	stop := s.Update(Input{Kind: Next}, now)
	if len(stop) != 1 || stop[0].Kind != StopOutput || len(s.Attempt.Marks) != 0 {
		t.Fatal("navigation did not clear and cancel")
	}
	if got := s.Update(Input{Kind: Done, ID: fx.ID}, now); len(got) != 0 {
		t.Fatal("old completion escaped stop barrier")
	}
	next := s.Update(Input{Kind: Stopped, ID: fx.ID}, now)
	if len(next) != 1 || next[0].Kind != PlayCW || s.Line != 1 {
		t.Fatal("new reference did not start after cleanup")
	}
}
func TestDisconnectAndPauseInvalidateCapture(t *testing.T) {
	for _, action := range []Kind{InputLost, TogglePause} {
		s := New(1)
		s.SendBack = true
		s.Phase = Ready
		now := time.Now()
		s.Update(Input{Kind: BeginCapture, At: time.Second}, now)
		s.Update(Input{Kind: KeyEvent, Event: keying.Event{At: time.Second, Down: true}}, now)
		s.Update(Input{Kind: action}, now)
		if s.Attempt.Active || s.Attempt.Err == nil {
			t.Fatal("capture survived interruption")
		}
		if got := s.Update(Input{Kind: ReplayAttempt}, now); len(got) != 0 {
			t.Fatal("incomplete attempt replayed")
		}
	}
}
