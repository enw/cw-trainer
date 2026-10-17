package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/session"
	"fmt"
	"strings"
)

type operatingReference struct{ page, scroll int }

func (m *Model) referenceKey(key string) tea.Cmd {
	r := m.reference
	switch key {
	case "space":
		if m.screen == practice {
			return m.dispatch(session.Input{Kind: session.TogglePause})
		}
	case "esc", "o":
		m.reference = nil
	case "j", "down":
		r.scroll = min(r.scroll+1, max(0, len(m.referenceRows())-max(1, m.height-7)))
	case "k", "up":
		r.scroll = max(0, r.scroll-1)
	case "i":
		if m.screen == practice {
			m.notice = "Esc closes reference; Esc returns to browsing before import"
			return nil
		}
		return m.importLists(library.PhoneticStarter())
	}
	return nil
}

func (m *Model) referenceRows() []string {
	rows := []string{"Each letter keeps its CW sound; speech uses its phonetic word."}
	for i, word := range library.PhoneticWords() {
		rows = append(rows, fmt.Sprintf("%c  %s", 'A'+i, word))
	}
	return rows
}

func (m *Model) referenceView() string {
	rows := m.referenceRows()
	height := max(1, m.height-7)
	start := min(m.reference.scroll, max(0, len(rows)-height))
	end := min(len(rows), start+height)
	return fmt.Sprintf("Operating reference · NATO / ITU alphabet\nSource: ARRL Quick Reference Operating Aids\n\n%s\n\nj/k or arrows scroll · i import this list (browse only)\no/Esc return · Space pause/resume practice · q quit", strings.Join(rows[start:end], "\n"))
}
