package session

import "math/rand/v2"

// SetExchanges distinguishes manual multi-line practice from ordinary items.
// Call only while configuring a stopped session.
func (s *Session) SetExchanges(exchanges []bool) {
	if s.Phase == Idle || s.Phase == Failed {
		s.exchangeEntries = append([]bool(nil), exchanges...)
	}
}
func (s *Session) resetRound() { s.shuffleBag = nil; s.roundStarted = false }
func (s *Session) isExchange() bool {
	return s.Entry < len(s.exchangeEntries) && s.exchangeEntries[s.Entry]
}
func (s *Session) advanceUnit() {
	if s.isExchange() {
		if s.Settings.LineAdvancement == "manual" {
			return
		}
		if s.Line+1 < s.lines {
			s.Line++
			s.Count = 0
			s.Revealed = false
			return
		}
		s.Line = 0
		s.Count = 0
		s.Revealed = false
	}
	s.advanceEntry()
}
func (s *Session) advanceEntry() {
	n := len(s.entryLengths)
	if n == 0 || s.Settings.Sequence == "repeat" {
		return
	}
	next := (s.Entry + 1) % n
	if s.Settings.Sequence == "shuffle" {
		if len(s.shuffleBag) == 0 {
			for i := 0; i < n; i++ {
				if s.roundStarted || i != s.Entry {
					s.shuffleBag = append(s.shuffleBag, i)
				}
			}
			if s.Shuffle != nil {
				s.Shuffle(s.shuffleBag)
			} else {
				rand.Shuffle(len(s.shuffleBag), func(i, j int) { s.shuffleBag[i], s.shuffleBag[j] = s.shuffleBag[j], s.shuffleBag[i] })
			}
			s.roundStarted = true
		}
		if len(s.shuffleBag) == 0 {
			next = s.Entry
		} else {
			next = s.shuffleBag[0]
			s.shuffleBag = s.shuffleBag[1:]
		}
	}
	if next != s.Entry || s.Line != 0 {
		s.Revealed = false
	}
	s.Entry = next
	s.Line = 0
	s.lines = s.entryLengths[next]
	s.Count = 0
}
