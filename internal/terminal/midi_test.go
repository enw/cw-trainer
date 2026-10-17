package terminal

import (
	"cw-looper/internal/library"
	"cw-looper/internal/midi"
	"cw-looper/internal/output"
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
