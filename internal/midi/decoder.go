// Package midi provides timestamped key events, independent of terminal focus.
package midi

import (
	"cw-looper/internal/keying"
	"time"
)

type Event = keying.Event

type Device struct {
	ID   int32
	Name string
}
type Port interface {
	Poll() ([]Event, error)
	Close() error
}
type Driver interface {
	Devices() ([]Device, error)
	Connect(int32) (Port, error)
}

// Decoder accepts channel-one notes 0/1/2 from Vail firmware variants. Neither
// note number nor the user's keyer setting substitutes for measured duration.
type Decoder struct {
	status byte
	data   [2]byte
	used   int
	sysex  bool
}

func (d *Decoder) Feed(at time.Duration, bytes []byte) []Event {
	var events []Event
	for _, b := range bytes {
		if b >= 0xf8 {
			continue
		}
		if b&0x80 != 0 {
			d.used = 0
			if b >= 0xf0 {
				d.status = 0
				d.sysex = b == 0xf0
				continue
			}
			d.status = b
			d.sysex = false
			continue
		}
		if d.sysex || d.status == 0 {
			continue
		}
		d.data[d.used] = b
		d.used++
		needed := 2
		if d.status&0xf0 == 0xc0 || d.status&0xf0 == 0xd0 {
			needed = 1
		}
		if d.used < needed {
			continue
		}
		if (d.status == 0x90 || d.status == 0x80) && d.data[0] <= 2 {
			events = append(events, Event{At: at, Note: d.data[0], Down: d.status == 0x90 && d.data[1] != 0})
		}
		d.used = 0
	}
	return events
}
