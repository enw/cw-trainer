package midi_test

import (
	"cw-looper/internal/midi"
	"errors"
	"testing"
)

type port struct {
	closed bool
	fail   bool
}

func (p *port) Poll() ([]midi.Event, error) {
	if p.fail {
		return nil, errors.New("removed")
	}
	return []midi.Event{{Down: true, Note: 1}}, nil
}
func (p *port) Close() error { p.closed = true; return nil }

type driver struct{ port *port }

func (d driver) Devices() ([]midi.Device, error)  { return []midi.Device{{ID: 1, Name: "fixture"}}, nil }
func (d driver) Connect(int32) (midi.Port, error) { return d.port, nil }
func TestDisconnectInvalidatesInputAndClosingPreventsLateConnect(t *testing.T) {
	p := &port{}
	m := midi.NewManager(driver{p})
	if err := m.Connect(1); err != nil {
		t.Fatal(err)
	}
	p.fail = true
	if _, err := m.Poll(); err == nil || !p.closed {
		t.Fatal("disconnect did not close port")
	}
	m.Close()
	if err := m.Connect(1); err == nil {
		t.Fatal("late connection stranded after close")
	}
}
