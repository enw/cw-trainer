package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"strings"
)

func (m *Model) lineForm(add bool) tea.Cmd {
	line := library.Line{Station: "learner"}
	if !add {
		line = m.currentLines()[m.cursor]
	}
	return m.openForm("Exchange line · review custom speech after CW edits",
		[]string{"CW text", "Station: learner / other", "Spoken override (clear to use generated speech)"},
		[]string{line.CW, line.Station, line.Spoken}, func(v []string) (library.Library, error) {
			data := clone(m.data)
			updated := line
			if add {
				updated.ID = id()
			}
			updated.CW = strings.ToUpper(strings.TrimSpace(v[0]))
			updated.Station = v[1]
			updated.Spoken = v[2]
			entry := &data.Lists[m.list].Entries[m.entry]
			if add {
				entry.Lines = append(entry.Lines, updated)
			} else {
				entry.Lines[m.cursor] = updated
			}
			return data, nil
		})
}

func (m *Model) exchangeForm(add bool) tea.Cmd {
	if !add {
		title := m.data.Lists[m.list].Entries[m.cursor].Title
		return m.openForm("Exchange title · Enter opens its line editor", []string{"Title"}, []string{title}, func(v []string) (library.Library, error) {
			data := clone(m.data)
			data.Lists[m.list].Entries[m.cursor].Title = v[0]
			return data, nil
		})
	}
	return m.openForm("New exchange · add more lines after saving",
		[]string{"Title", "First line CW text", "Station: learner / other", "Spoken override (clear to use generated speech)"},
		[]string{"", "", "learner", ""}, func(v []string) (library.Library, error) {
			data := clone(m.data)
			entry := library.Entry{ID: id(), Kind: "exchange", Title: v[0], Lines: []library.Line{{ID: id(), Station: v[2], CW: strings.ToUpper(strings.TrimSpace(v[1])), Spoken: v[3]}}}
			data.Lists[m.list].Entries = append(data.Lists[m.list].Entries, entry)
			return data, nil
		})
}
