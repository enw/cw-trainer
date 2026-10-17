package session_test

import (
	"cw-looper/internal/session"
	"errors"
	"testing"
	"time"
)

func TestNavigationWaitsForStopAndRejectsStaleCompletion(t *testing.T) {
	s := session.New(2)
	now := time.Unix(0, 0)
	first := s.Update(session.Input{Kind: session.Start}, now)[0]
	stop := s.Update(session.Input{Kind: session.Next}, now)
	if len(stop) != 1 || stop[0].Kind != session.StopOutput {
		t.Fatal("navigation did not cancel output")
	}
	if s.Line != 1 || s.Phase != session.Stopping {
		t.Fatal("wrong target or cancellation state")
	}
	if got := s.Update(session.Input{Kind: session.Done, ID: first.ID}, now); len(got) != 0 {
		t.Fatal("stale done launched output")
	}
	s.Update(session.Input{Kind: session.Previous}, now)
	next := s.Update(session.Input{Kind: session.Stopped, ID: first.ID}, now)
	if len(next) != 1 || next[0].Kind != session.PlayCW || s.Line != 0 {
		t.Fatal("did not coalesce navigation")
	}
	if next[0].ID == first.ID {
		t.Fatal("generation not changed")
	}
}

func TestPauseRestartsPortionAndFreezesGap(t *testing.T) {
	s := session.New(2)
	now := time.Unix(0, 0)
	first := s.Update(session.Input{Kind: session.Start}, now)[0]
	s.Update(session.Input{Kind: session.TogglePause}, now.Add(time.Second))
	s.Update(session.Input{Kind: session.Stopped, ID: first.ID}, now.Add(time.Second))
	if s.Phase != session.Paused {
		t.Fatal("not paused")
	}
	restarted := s.Update(session.Input{Kind: session.TogglePause}, now.Add(2*time.Second))[0]
	if restarted.Kind != session.PlayCW || s.Count != 0 {
		t.Fatal("partial CW counted")
	}
	gap := s.Update(session.Input{Kind: session.Done, ID: restarted.ID}, now.Add(3*time.Second))[0]
	if s.Count != 1 || gap.Duration != 2*time.Second {
		t.Fatal("wrong repetition gap")
	}
	s.Update(session.Input{Kind: session.TogglePause}, now.Add(3500*time.Millisecond))
	s.Update(session.Input{Kind: session.Stopped, ID: gap.ID}, now.Add(4*time.Second))
	resumed := s.Update(session.Input{Kind: session.TogglePause}, now.Add(10*time.Second))[0]
	if resumed.Kind != session.Wait || resumed.Duration != 1500*time.Millisecond {
		t.Fatalf("gap not frozen: %+v", resumed)
	}
}

func TestSpeechAndNavigationWhilePaused(t *testing.T) {
	s := session.New(2)
	now := time.Now()
	cw := s.Update(session.Input{Kind: session.Start}, now)[0]
	s.Update(session.Input{Kind: session.TogglePause}, now)
	s.Update(session.Input{Kind: session.Stopped, ID: cw.ID}, now)
	speech := s.Update(session.Input{Kind: session.Speak}, now)[0]
	if speech.Kind != session.PlaySpeech {
		t.Fatal("no speech")
	}
	if fx := s.Update(session.Input{Kind: session.Done, ID: speech.ID}, now); len(fx) != 0 || s.Phase != session.Paused {
		t.Fatal("paused speech auto-resumed CW")
	}
	if fx := s.Update(session.Input{Kind: session.Next}, now); len(fx) != 0 || s.Line != 1 {
		t.Fatal("paused navigation played")
	}
	if fx := s.Update(session.Input{Kind: session.TogglePause}, now); fx[0].Kind != session.PlayCW {
		t.Fatal("did not resume new line")
	}
}

func TestQuitWinsDuringCancellation(t *testing.T) {
	s := session.New(2)
	now := time.Now()
	cw := s.Update(session.Input{Kind: session.Start}, now)[0]
	s.Update(session.Input{Kind: session.Speak}, now)
	s.Update(session.Input{Kind: session.Quit}, now)
	s.Update(session.Input{Kind: session.Next}, now)
	result := s.Update(session.Input{Kind: session.Stopped, ID: cw.ID}, now)
	if len(result) != 1 || result[0].Kind != session.Exit {
		t.Fatal("quit lost to pending playback")
	}
}

func TestSpeechInterruptionAndOutputErrors(t *testing.T) {
	s := session.New(2)
	now := time.Now()
	cw := s.Update(session.Input{Kind: session.Start}, now)[0]
	s.Update(session.Input{Kind: session.Speak}, now)
	speech := s.Update(session.Input{Kind: session.Stopped, ID: cw.ID}, now)[0]
	if speech.Kind != session.PlaySpeech {
		t.Fatal("wrong transition")
	}
	s.Update(session.Input{Kind: session.Next}, now)
	s.Update(session.Input{Kind: session.Done, ID: speech.ID}, now)
	next := s.Update(session.Input{Kind: session.Stopped, ID: speech.ID}, now)[0]
	if next.Kind != session.PlayCW || s.Line != 1 {
		t.Fatal("speech navigation failed")
	}
	fx := s.Update(session.Input{Kind: session.Error, ID: next.ID, Err: errors.New("output failed")}, now)
	if len(fx) != 0 || s.Phase != session.Failed || s.Err == nil {
		t.Fatal("failure loops or disappears")
	}
}
