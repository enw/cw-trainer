package terminal

import (
	"cw-looper/internal/library"
	"cw-looper/internal/midi"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"errors"
	"testing"
	"time"
)

func TestMIDIDiagnosticRejectsOldEventsAndClearsDisconnect(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.midi.connected = true
	m.midi.generation = 2
	m.midi.down = map[byte]time.Duration{}
	m.midiMessage(midiPollMsg{generation: 1, events: []midi.Event{{Note: 1, Down: true}}})
	if m.midi.received != 0 {
		t.Fatal("stale connection event accepted")
	}
	m.midiMessage(midiPollMsg{generation: 2, events: []midi.Event{{Note: 1, Down: true}}})
	if len(m.midi.down) != 1 {
		t.Fatal("live mark missing")
	}
	m.midiMessage(midiPollMsg{generation: 2, err: errors.New("unplugged")})
	if m.midi.connected || len(m.midi.down) != 0 {
		t.Fatal("held note survived disconnect")
	}
}

func TestLateMIDIAfterFinishDeliversReplayStopAcknowledgment(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.session.SendBack = true
	m.session.Phase = session.Ready
	now := time.Now()
	for _, input := range []session.Input{
		{Kind: session.BeginCapture, At: time.Second},
		{Kind: session.KeyEvent, Event: midi.Event{At: time.Second, Down: true}},
		{Kind: session.KeyEvent, Event: midi.Event{At: 1080 * time.Millisecond}},
		{Kind: session.FinishCapture, At: 2 * time.Second},
		{Kind: session.CaptureDrained, At: 2 * time.Second},
		{Kind: session.ReplayAttempt},
	} {
		m.session.Update(input, now)
	}
	cmd := m.midiEvent(midi.Event{At: 1500 * time.Millisecond, Down: true})
	if cmd == nil || m.session.Phase != session.Stopping {
		t.Fatal("late input did not deliver replay cancellation command")
	}
}
