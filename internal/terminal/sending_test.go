package terminal

import (
	"cw-looper/internal/library"
	"cw-looper/internal/midi"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"strings"
	"testing"
	"time"
)

func TestSendingFeedbackRespectsHiddenTargetAndIndependentSpeed(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.session.SendBack = true
	m.session.Phase = session.Captured
	m.session.Settings.Display = "hidden"
	m.session.Attempt.Begin(time.Second)
	m.session.Attempt.Feed(midi.Event{At: time.Second, Down: true})
	m.session.Attempt.Feed(midi.Event{At: 1080 * time.Millisecond})
	_ = m.session.Attempt.Finish()
	hidden := m.sendingView()
	if !strings.Contains(hidden, "Tentative decode: E") || strings.Contains(hidden, "Different Morse") || strings.Contains(hidden, "Matching Morse") {
		t.Fatal(hidden)
	}
	m.session.Revealed = true
	if !strings.Contains(m.sendingView(), "Different Morse") {
		t.Fatal("reveal did not compare")
	}
	before := m.session.Settings.CharacterWPM
	m.sendingKey("+")
	if m.sendWPM != 16 || m.session.Settings.CharacterWPM != before {
		t.Fatal("sending speed coupled to reference")
	}
}

func TestFullFeedbackCanScrollThroughLongAttemptWithoutRevealingTarget(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.session.SendBack = true
	m.session.Phase = session.Captured
	m.session.Settings.Display = "hidden"
	at := time.Second
	m.session.Attempt.Begin(at)
	for i := 0; i < 1000; i++ {
		m.session.Attempt.Feed(midi.Event{At: at, Down: true})
		m.session.Attempt.Feed(midi.Event{At: at + 20*time.Millisecond})
		at += 80 * time.Millisecond
	}
	_ = m.session.Attempt.Finish()
	m.sendWPM = 60
	m.width = 80
	m.height = 24
	m.sendingKey("f")
	all := strings.Join(m.feedbackRows(), "\n")
	if strings.Contains(all, "Target: ") || strings.Contains(all, "Different Morse") {
		t.Fatal("hidden target leaked")
	}
	if strings.Count(all, "E") < 1000 {
		t.Fatal("long decode truncated")
	}
	// At the supported minimum screen size, footer and data stay within bounds.
	if len(strings.Split(m.feedbackView(), "\n")) > 24 {
		t.Fatal("feedback exceeds screen")
	}
	for i := 0; i < 100; i++ {
		m.feedbackKey("j")
	}
	if m.feedbackScroll == 0 || !strings.Contains(m.feedbackView(), "Gap ideals:") {
		t.Fatal("cannot reach feedback end")
	}
	m.feedbackKey("esc")
	if m.feedbackOpen {
		t.Fatal("feedback did not close")
	}
}

func TestFeedbackDoesNotBlockActiveCaptureOrReplayPause(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.session.SendBack = true
	for _, phase := range []session.Phase{session.CW, session.Ready, session.Capturing} {
		m.session.Phase = phase
		m.sendingKey("f")
		if m.feedbackOpen {
			t.Fatalf("feedback opened during %s", phase)
		}
	}
	m.session.Phase = session.Replay
	m.sendingKey("f")
	if !m.feedbackOpen {
		t.Fatal("replay feedback unavailable")
	}
	cmd := m.feedbackKey("space")
	if cmd == nil {
		t.Fatal("pause acknowledgment missing")
	}
	m.Update(cmd())
	if m.session.Phase != session.Paused {
		t.Fatalf("pause ignored: %s", m.session.Phase)
	}
}
