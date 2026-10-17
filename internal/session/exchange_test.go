package session

import (
	"testing"
	"time"
)

func TestAutomaticExchangeRetainsLineOrderAndWraps(t *testing.T) {
	for _, sequence := range []string{"repeat", "cycle", "shuffle"} {
		t.Run(sequence, func(t *testing.T) {
			s := New(2)
			s.SelectEntries([]int{2, 1}, 0, 0)
			s.SetExchanges([]bool{true, false})
			s.Settings.LineAdvancement = "automatic"
			s.Settings.Sequence = sequence
			now := time.Now()
			fx := s.Update(Input{Kind: Start}, now)[0]
			advance := func() {
				fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
				fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
			}
			s.Revealed = true
			advance()
			if s.Entry != 0 || s.Line != 1 || s.Revealed {
				t.Fatal("line order or reveal reset incorrect")
			}
			advance()
			want := 1
			if sequence == "repeat" {
				want = 0
			}
			if s.Entry != want || s.Line != 0 {
				t.Fatal("incorrect exchange boundary")
			}
		})
	}
}
func TestManualExchangeGatesTopLevelCycle(t *testing.T) {
	s := New(2)
	s.SelectEntries([]int{2, 1}, 0, 1)
	s.SetExchanges([]bool{true, false})
	s.Settings.Sequence = "cycle"
	now := time.Now()
	fx := s.Update(Input{Kind: Start}, now)[0]
	for i := 0; i < 4; i++ {
		fx = s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
	}
	if s.Entry != 0 || s.Line != 1 {
		t.Fatal("manual exchange advanced")
	}
}

func TestAllStyleSequenceAndLineModeCombinations(t *testing.T) {
	for _, style := range []string{"cw", "sandwich", "occasional"} {
		for _, sequence := range []string{"repeat", "cycle", "shuffle"} {
			for _, lineMode := range []string{"manual", "automatic"} {
				t.Run(style+"/"+sequence+"/"+lineMode, func(t *testing.T) {
					s := New(2)
					s.SelectEntries([]int{2, 1}, 0, 0)
					s.SetExchanges([]bool{true, false})
					s.Settings.Style = style
					s.Settings.Sequence = sequence
					s.Settings.LineAdvancement = lineMode
					s.Settings.SpeechInterval = 2
					now := time.Now()
					fx := s.Update(Input{Kind: Start}, now)[0]
					portions := 1
					if style != "cw" {
						portions = 3
					}
					for unit := 0; unit < 2; unit++ {
						for portion := 0; portion < portions; portion++ {
							if fx.Kind != PlayCW && fx.Kind != PlaySpeech {
								t.Fatal("missing portion")
							}
							gap := s.Update(Input{Kind: Done, ID: fx.ID}, now)[0]
							if gap.Kind != Wait {
								t.Fatal("missing boundary gap")
							}
							fx = s.Update(Input{Kind: Done, ID: gap.ID}, now)[0]
						}
						if unit == 0 && lineMode == "automatic" && s.Line != 1 {
							t.Fatal("did not advance ordered line")
						}
					}
					if lineMode == "manual" {
						if s.Entry != 0 || s.Line != 0 {
							t.Fatal("manual advanced")
						}
					} else if sequence == "repeat" {
						if s.Entry != 0 || s.Line != 0 {
							t.Fatal("repeat failed")
						}
					} else if s.Entry != 1 {
						t.Fatal("did not sequence after exchange")
					}
				})
			}
		}
	}
}
