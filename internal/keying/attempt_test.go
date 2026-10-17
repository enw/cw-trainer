package keying_test

import (
	"cw-looper/internal/keying"
	"testing"
	"time"
)

func TestCapturedTimingReplaysWithoutQuantization(t *testing.T) {
	var attempt keying.Attempt
	attempt.Begin(time.Second)
	for _, e := range []keying.Event{{At: 2 * time.Second, Note: 1, Down: true}, {At: 2081 * time.Millisecond, Note: 1}, {At: 2174 * time.Millisecond, Note: 2, Down: true}, {At: 2413 * time.Millisecond, Note: 2}} {
		attempt.Feed(e)
	}
	if err := attempt.Finish(); err != nil {
		t.Fatal(err)
	}
	timeline, err := attempt.Timeline()
	if err != nil {
		t.Fatal(err)
	}
	if timeline.Samples != 19824 || len(timeline.Segments) != 3 {
		t.Fatalf("actual timing altered: %+v", timeline)
	}
	if timeline.Segments[0].End != 3888 || timeline.Segments[1].Tone || timeline.Segments[1].End != 8352 {
		t.Fatal("mark or gap quantized")
	}
}

func TestIncompleteAndBoundedCapturesCannotReplay(t *testing.T) {
	for _, kind := range []string{"empty", "held", "overlap", "unmatched", "reversed", "timeout", "limit"} {
		t.Run(kind, func(t *testing.T) {
			var a keying.Attempt
			a.Begin(time.Second)
			switch kind {
			case "held":
				a.Feed(keying.Event{At: time.Second, Down: true})
			case "overlap":
				a.Feed(keying.Event{At: time.Second, Down: true})
				a.Feed(keying.Event{At: 2 * time.Second, Note: 1, Down: true})
			case "unmatched":
				a.Feed(keying.Event{At: time.Second})
			case "reversed":
				a.Feed(keying.Event{At: 2 * time.Second, Down: true})
				a.Feed(keying.Event{At: time.Second})
			case "timeout":
				a.Check(122 * time.Second)
			case "limit":
				for i := 0; i < 1025; i++ {
					at := time.Second + time.Duration(i)*2*time.Millisecond
					a.Feed(keying.Event{At: at, Down: true})
					a.Feed(keying.Event{At: at + time.Millisecond})
				}
			}
			if a.Finish() == nil {
				t.Fatal("invalid capture accepted")
			}
			if _, err := a.Timeline(); err == nil {
				t.Fatal("incomplete capture replayable")
			}
		})
	}
}
