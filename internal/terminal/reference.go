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
	case "tab", "h", "l", "left", "right":
		r.page = 1 - r.page
		r.scroll = 0
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
		if r.page == 1 {
			return m.importLists(library.RSTStarter())
		}
		return m.importLists(library.PhoneticStarter())
	}
	return nil
}

func (m *Model) referenceRows() []string {
	if m.reference.page == 1 {
		rows := []string{"Order: Readability, Strength, Tone. CW: three ratings; voice: R/S.", "Clean practice tones teach meanings, not simulated signal quality."}
		for _, scale := range library.RSTScales() {
			rows = append(rows, scale.Name)
			for i, meaning := range scale.Meanings {
				rows = append(rows, fmt.Sprintf("%d  %s", i+1, meaning))
			}
		}
		return rows
	}
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
	title := "NATO / ITU alphabet"
	if m.reference.page == 1 {
		title = "RST · Readability / Strength / Tone"
	}
	return fmt.Sprintf("Operating reference · %s\nSource: ARRL Quick Reference Operating Aids\n\n%s\n\nj/k or arrows scroll · i import this list (browse only)\nh/l or ←/→ tabs · o/Esc return · Space pause/resume · q quit", title, strings.Join(rows[start:end], "\n"))
}
