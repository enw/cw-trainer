package keying

import (
	"strings"
	"testing"
	"time"
)

func measuredAttempt(durations ...int) Attempt {
	var a Attempt
	at := time.Second
	a.Begin(at)
	for i, ms := range durations {
		end := at + time.Duration(ms)*time.Millisecond
		if i%2 == 0 {
			a.Feed(Event{At: at, Down: true})
			a.Feed(Event{At: end})
		}
		at = end
	}
	_ = a.Finish()
	return a
}

func TestFeedbackDecodesMeasuredSendingIndependently(t *testing.T) {
	// A E at 15 WPM: dit, dah, word gap, dit.
	a := measuredAttempt(80, 80, 240, 560, 80)
	f := a.Feedback(15)
	if f.Text != "A E" || f.Compare("A E") != "Matching Morse patterns and word boundaries" {
		t.Fatalf("feedback: %+v", f)
	}
	if f.Compare("A T") != "Different Morse patterns or word boundaries" {
		t.Fatal("expected answer influenced decoding")
	}
}

func TestFeedbackUncertaintySpeedAndContinuousProsigns(t *testing.T) {
	cases := []struct {
		name      string
		durations []int
		wpm       float64
		text      string
		uncertain bool
	}{
		{"uncertain gap", []int{80, 160, 80}, 15, "E[gap?]E", true},
		{"uncertain mark", []int{160}, 15, "[??]", true},
		{"20 WPM A", []int{60, 60, 180}, 20, "A", false},
		{"wrong speed", []int{60, 60, 180}, 15, "[.??]", true},
		{"continuous AR", []int{80, 80, 240, 80, 80, 80, 240, 80, 80}, 15, "<AR>/+", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := measuredAttempt(tc.durations...).Feedback(tc.wpm)
			if f.Text != tc.text || f.Uncertain != tc.uncertain {
				t.Fatalf("%+v", f)
			}
			if tc.uncertain && f.Compare("E") == "Matching Morse patterns and word boundaries" {
				t.Fatal("uncertainty became match")
			}
		})
	}
	ar := measuredAttempt(80, 80, 240, 80, 80, 80, 240, 80, 80).Feedback(15)
	if ar.Compare("<AR>") != "Matching Morse patterns and word boundaries" || ar.Compare("AR") != "Different Morse patterns or word boundaries" {
		t.Fatal("prosign boundary lost")
	}
	for _, target := range []string{"E", "EEE", "E E", "T"} {
		f := measuredAttempt(80, 240, 80).Feedback(15)
		if f.Compare(target) != "Different Morse patterns or word boundaries" {
			t.Fatalf("missed extra/missing/wrong symbol or gap: %s", target)
		}
	}
}

func TestFeedbackSeparatesNoInputAndIncompleteCapture(t *testing.T) {
	var a Attempt
	if a.Feedback(15).Status != "No input" {
		t.Fatal("new attempt")
	}
	a.Begin(time.Second)
	_ = a.Finish()
	if a.Feedback(15).Status != "No input" {
		t.Fatal("empty attempt")
	}
	a.Begin(time.Second)
	a.Feed(Event{At: time.Second, Down: true})
	_ = a.Finish()
	if !strings.HasPrefix(a.Feedback(15).Status, "Incomplete capture:") {
		t.Fatal("held input")
	}
	a.Invalidate("MIDI disconnected")
	if !strings.HasPrefix(a.Feedback(15).Status, "Incomplete capture:") {
		t.Fatal("disconnect")
	}
}

func TestFeedbackAcceptsObservedAdapterDurations(t *testing.T) {
	// Representative recorded 2026-10-16 dit and dah durations; gap is synthetic.
	var a Attempt
	a.Begin(time.Second)
	for _, e := range []Event{
		{At: time.Second, Down: true}, {At: time.Second + 79795*time.Microsecond},
		{At: time.Second + 159795*time.Microsecond, Down: true},
		{At: time.Second + 399434*time.Microsecond},
	} {
		a.Feed(e)
	}
	_ = a.Finish()
	f := a.Feedback(15)
	if f.Text != "A" || f.Uncertain {
		t.Fatalf("observed lengths rejected: %+v", f)
	}
}
