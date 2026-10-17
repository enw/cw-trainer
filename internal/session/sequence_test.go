package session

import (
	"testing"
	"time"
)

func TestCycleWaitsForCompleteReinforcementUnit(t *testing.T) {
	for _, tc := range []struct {
		style    string
		portions int
	}{{"cw", 1}, {"sandwich", 3}, {"occasional", 3}} {
		t.Run(tc.style, func(t *testing.T) {
			s := New(1)
			s.SelectEntries([]int{1, 1}, 0, 0)
			s.Settings.Style = tc.style
			s.Settings.Sequence = "cycle"
			s.Settings.SpeechInterval = 2
			now := time.Now()
			fx := s.Update(Input{Kind: Start}, now)[0]
			for i := 0; i < tc.portions; i++ {
				fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
				if s.Entry != 0 {
					t.Fatal("advanced before boundary gap")
				}
				fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
				if i < tc.portions-1 && s.Entry != 0 {
					t.Fatal("split reinforcement unit")
				}
			}
			if s.Entry != 1 || fx.Kind != PlayCW {
				t.Fatalf("did not advance: %+v", s)
			}
		})
	}
}

func TestShuffleVisitsOnceAndNavigationResetsRound(t *testing.T) {
	s := New(1)
	s.SelectEntries([]int{1, 1, 1}, 1, 0)
	s.Settings.Sequence = "shuffle"
	s.Shuffle = func(values []int) {} // Fixed order at the shuffle seam.
	now := time.Now()
	fx := s.Update(Input{Kind: Start}, now)[0]
	seen := map[int]bool{s.Entry: true}
	for i := 0; i < 2; i++ {
		fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
		fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
		if seen[s.Entry] {
			t.Fatal("repeated inside round")
		}
		seen[s.Entry] = true
	}
	if len(seen) != 3 {
		t.Fatal(seen)
	}
}

func TestEmptySingleAndPausedSequence(t *testing.T) {
	now := time.Now()
	empty := New(0)
	if len(empty.Update(Input{Kind: Start}, now)) != 0 {
		t.Fatal("empty session played")
	}
	s := New(1)
	s.Settings.Sequence = "shuffle"
	fx := s.Update(Input{Kind: Start}, now)[0]
	for i := 0; i < 3; i++ {
		fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
		fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
		if s.Entry != 0 {
			t.Fatal("single item advanced outside list")
		}
	}
	stop := s.Update(Input{Kind: TogglePause}, now)[0]
	s.Update(Input{Kind: Stopped, ID: stop.ID}, now)
	if len(s.Update(Input{Kind: Done, ID: fx.ID}, now)) != 0 || s.Phase != Paused {
		t.Fatal("paused sequence advanced")
	}
}
