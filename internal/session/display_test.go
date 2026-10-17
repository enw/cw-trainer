package session

import (
	"testing"
	"time"
)

func TestRevealAndProgressBelongToCurrentLine(t *testing.T) {
	s := New(2)
	now := time.Now()
	e := s.Update(Input{Kind: Start}, now)[0]
	s.Update(Input{Kind: Reveal}, now)
	s.Update(Input{Kind: Progress, ID: e.ID, Sample: 123}, now)
	if !s.Revealed || s.Sample != 123 {
		t.Fatal("active display events lost")
	}
	s.Update(Input{Kind: Next}, now)
	s.Update(Input{Kind: Progress, ID: e.ID, Sample: 456}, now)
	if s.Revealed || s.Sample != -1 {
		t.Fatal("stale display leaked after navigation")
	}
}
