package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
)

func (m *Model) importStarter() tea.Cmd {
	m.confirm = &confirmation{"Add a fresh copy of the four starter lists? Existing lists are preserved.", func() tea.Cmd {
		data := clone(m.data)
		starter := library.Starter()
		for _, list := range starter.Lists {
			list.ID = id()
			for i := range list.Entries {
				list.Entries[i].ID = id()
				for j := range list.Entries[i].Lines {
					list.Entries[i].Lines[j].ID = id()
				}
			}
			data.Lists = append(data.Lists, list)
		}
		if err := data.Validate(); err != nil {
			m.notice = err.Error()
			return nil
		}
		m.data = data
		m.dirty = true
		return m.persist()
	}}
	return nil
}
