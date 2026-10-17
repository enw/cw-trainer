package session

import (
	"testing"
	"time"
)

func TestReinforcementOrder(t *testing.T) {
	for _, tc := range []struct {
		style string
		want  []EffectKind
	}{
		{"cw", []EffectKind{PlayCW, Wait, PlayCW, Wait, PlayCW}},
		{"sandwich", []EffectKind{PlayCW, Wait, PlaySpeech, Wait, PlayCW, Wait, PlayCW}},
		{"occasional", []EffectKind{PlayCW, Wait, PlayCW, Wait, PlaySpeech, Wait, PlayCW}},
	} {
		t.Run(tc.style, func(t *testing.T) {
			s := New(1)
			s.Settings.Style = tc.style
			s.Settings.SpeechInterval = 2
			now := time.Now()
			fx := s.Update(Input{Kind: Start}, now)
			for i, want := range tc.want {
				if len(fx) != 1 || fx[0].Kind != want {
					t.Fatalf("portion %d: got %+v want %v", i, fx, want)
				}
				fx = s.Update(Input{Kind: Done, ID: fx[0].ID}, now)
			}
		})
	}
}

func TestIntervalWaitsForNextUnitAndPolicyChangeStops(t *testing.T) {
	s := New(1)
	s.Settings.Style = "occasional"
	s.Settings.SpeechInterval = 2
	now := time.Now()
	fx := s.Update(Input{Kind: Start}, now)[0]
	config := s.Settings
	config.SpeechInterval = 1
	if len(s.Configure(config, now)) != 0 {
		t.Fatal("interval interrupted CW")
	}
	fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
	fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
	if fx.Kind != PlayCW {
		t.Fatal("new interval changed active unit")
	}
	config.Style = "sandwich"
	stop := s.Configure(config, now)
	if len(stop) != 1 || stop[0].Kind != StopOutput {
		t.Fatal("mode change bypassed cancellation")
	}
	if got := s.Update(Input{Kind: Done, ID: fx.ID}, now); len(got) != 0 {
		t.Fatal("old completion launched output")
	}
	fx = s.Update(Input{Kind: Stopped, ID: fx.ID}, now)[0]
	if fx.Kind != PlayCW || s.Count != 0 {
		t.Fatal("new mode did not start fresh CW")
	}
}
