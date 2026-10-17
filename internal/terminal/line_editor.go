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
