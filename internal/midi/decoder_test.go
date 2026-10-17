package midi_test

import (
	"cw-looper/internal/midi"
	"testing"
	"time"
)

func TestRecordedVailEventsAndRunningStatus(t *testing.T) {
	var d midi.Decoder
	down := d.Feed(time.Millisecond*6949, []byte{0x90, 1, 127})
	up := d.Feed(time.Millisecond*7029, []byte{0x80, 1, 0})
	if len(down) != 1 || !down[0].Down || len(up) != 1 || up[0].Down || up[0].At-down[0].At != 80*time.Millisecond {
		t.Fatal("recorded dit lost")
	}
	got := d.Feed(time.Second, []byte{0x90, 2, 127, 0xf8, 2, 0, 0x91, 1, 127, 0xb0, 0, 0})
	if len(got) != 2 || !got[0].Down || got[1].Down {
		t.Fatalf("running status/release/channel filtering: %+v", got)
	}
}
