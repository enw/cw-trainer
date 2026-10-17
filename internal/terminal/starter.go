package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"fmt"
)

func (m *Model) importStarter() tea.Cmd {
	return m.importLists(library.Starter().Lists...)
}

func (m *Model) importLists(lists ...library.List) tea.Cmd {
	prompt := fmt.Sprintf("Add fresh copies of %d starter lists? Existing lists are preserved.", len(lists))
	if len(lists) == 1 {
		prompt = "Add a fresh copy of " + lists[0].Name + "? Existing lists are preserved."
	}
	m.confirm = &confirmation{prompt, func() tea.Cmd {
		data := clone(m.data)
		for _, list := range lists {
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
