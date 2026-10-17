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
