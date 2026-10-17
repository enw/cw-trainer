package morse_test

import (
	"bytes"
	"cw-looper/internal/morse"
	"encoding/binary"
	"io"
	"math"
	"testing"
)

func TestPARISCalibration(t *testing.T) {
	for _, tc := range []struct {
		overall float64
		samples int64
	}{{15, 192000}, {20, 144000}} {
		line, err := morse.Compile("PARIS", 20, tc.overall, true)
		if err != nil {
			t.Fatal(err)
		}
		if line.Samples != tc.samples {
			t.Fatalf("overall %g: got %d samples, want %d", tc.overall, line.Samples, tc.samples)
		}
	}
}

func TestRenderedPCM(t *testing.T) {
	timeline, err := morse.Compile("T", 20, 15, false)
	if err != nil {
		t.Fatal(err)
	}
	pcm, err := morse.NewPCM(timeline, 600)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(pcm)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 8640*4 {
		t.Fatalf("T should last 180 ms, got %d bytes", len(data))
	}
	crossings := 0
	previous := float32(0)
	for i := 0; i < len(data); i += 4 {
		sample := math.Float32frombits(binary.LittleEndian.Uint32(data[i:]))
		if sample > 0 && previous <= 0 {
			crossings++
		}
		previous = sample
		if math.Abs(float64(sample)) > 0.21 {
			t.Fatal("unbounded amplitude")
		}
	}
	if crossings < 107 || crossings > 109 {
		t.Fatalf("want 108 cycles at 600 Hz, got %d", crossings)
	}
	if _, err := pcm.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	again, _ := io.ReadAll(pcm)
	if !bytes.Equal(data, again) {
		t.Fatal("seek changed rendering")
	}
	if math.Float32frombits(binary.LittleEndian.Uint32(data[:4])) != 0 {
		t.Fatal("attack does not start at zero")
	}
}

func TestEncodingValidationAndGaps(t *testing.T) {
	for _, text := range []string{"", "CQ <XX>", "CQ!", "é", "A\tB"} {
		if _, err := morse.Compile(text, 20, 15, false); err == nil {
			t.Fatalf("accepted %q", text)
		}
	}
	timeline, _ := morse.Compile("ET", 20, 15, false)
	if timeline.Segments[0].End != 2880 {
		t.Fatal("E is not 60 ms")
	}
	if timeline.Segments[1].End-timeline.Segments[1].Start != 16219 {
		t.Fatal("wrong Farnsworth character gap")
	}
	timeline, _ = morse.Compile("E T", 20, 15, false)
	if timeline.Segments[1].End-timeline.Segments[1].Start != 37844 {
		t.Fatal("wrong Farnsworth word gap")
	}
}

func TestProsignIsContinuousAndDistinctFromLetters(t *testing.T) {
	joined, err := morse.Compile("<AR>", 20, 15, false)
	if err != nil {
		t.Fatal(err)
	}
	if joined.Samples != 37440 {
		t.Fatalf("AR joined should be 780 ms, got %d samples", joined.Samples)
	}
	letters, _ := morse.Compile("AR", 20, 15, false)
	if letters.Samples <= joined.Samples {
		t.Fatal("ordinary letters lost spacing")
	}
	mixed, err := morse.Compile("CQ DE KE2ISH <AR>", 20, 15, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(mixed.Markers) != 11 {
		t.Fatalf("prosign not a single display token: %d", len(mixed.Markers))
	}
}

func TestAllProsignDurationsHaveOnlyElementGaps(t *testing.T) {
	for token, units := range map[string]int64{"<AR>": 13, "<AS>": 11, "<BT>": 13, "<SK>": 15, "<HH>": 15} {
		timeline, err := morse.Compile(token, 20, 15, false)
		if err != nil {
			t.Fatal(err)
		}
		if timeline.Samples != units*2880 || len(timeline.Markers) != 1 {
			t.Fatalf("%s: wrong continuous timing", token)
		}
	}
}
