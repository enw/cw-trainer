package terminal

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"crypto/rand"
	"cw-looper/internal/library"
	"cw-looper/internal/morse"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type form struct {
	picker bool
	pick   int
	title  string
	labels []string
	fields []textinput.Model
	focus  int
	apply  func([]string) (library.Library, error)
}
type confirmation struct {
	prompt string
	yes    func() tea.Cmd
}
type savedMsg struct{ err error }

func clone(data library.Library) library.Library {
	encoded, _ := json.Marshal(data)
	var copy library.Library
	_ = json.Unmarshal(encoded, &copy)
	return copy
}
func id() string { return fmt.Sprintf("%x", rand.Text()) }
func (m *Model) openForm(title string, labels, values []string, apply func([]string) (library.Library, error)) tea.Cmd {
	f := &form{title: title, labels: labels, apply: apply}
	for _, v := range values {
		input := textinput.New()
		input.SetValue(v)
		input.CharLimit = 1024
		input.SetWidth(64)
		f.fields = append(f.fields, input)
	}
	m.form = f
	return f.fields[0].Focus()
}
func (m *Model) formKey(msg tea.KeyPressMsg) tea.Cmd {
	f := m.form
	if f.picker {
		return m.pickerKey(msg.String())
	}
	if msg.String() == "ctrl+p" && strings.Contains(f.labels[f.focus], "CW") {
		f.picker = true
		return nil
	}
	switch msg.String() {
	case "esc":
		m.form = nil
		return nil
	case "tab", "shift+tab":
		f.fields[f.focus].Blur()
		step := 1
		if msg.String() == "shift+tab" {
			step = -1
		}
		f.focus = (f.focus + step + len(f.fields)) % len(f.fields)
		return f.fields[f.focus].Focus()
	case "ctrl+s":
		values := make([]string, len(f.fields))
		for i := range f.fields {
			values[i] = f.fields[i].Value()
		}
		data, err := f.apply(values)
		if err == nil {
			err = data.Validate()
		}
		if err != nil {
			m.notice = err.Error()
			return nil
		}
		m.form = nil
		m.data = data
		var playback tea.Cmd
		if m.screen == practice {
			playback = m.applyEffects(m.session.Configure(m.data.Lists[m.list].Settings, time.Now()))
		}
		m.dirty = true
		return tea.Batch(playback, m.persist())
	}
	var cmd tea.Cmd
	f.fields[f.focus], cmd = f.fields[f.focus].Update(msg)
	return cmd
}
func (m *Model) formView() string {
	f := m.form
	if f.picker {
		return m.pickerView()
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", f.title)
	start := max(0, f.focus-3)
	for i := start; i < min(len(f.labels), start+5); i++ {
		label := f.labels[i]
		fmt.Fprintf(&b, "%s\n%s\n\n", label, f.fields[i].View())
	}
	b.WriteString("Tab next field · Ctrl+P prosigns · Ctrl+S save · Esc cancel")
	return b.String()
}
func (m *Model) persist() tea.Cmd {
	if m.save == nil {
		m.dirty = false
		return func() tea.Msg { return savedMsg{} }
	}
	m.saving = true
	data := clone(m.data)
	return func() tea.Msg { return savedMsg{m.save(data)} }
}
func (m *Model) confirmKey(key string) tea.Cmd {
	if key == "y" {
		action := m.confirm.yes
		m.confirm = nil
		return action()
	}
	if key == "esc" || key == "n" {
		m.confirm = nil
	}
	return nil
}
func (m *Model) editKey(key string) tea.Cmd {
	if m.screen == lines {
		if key == "e" {
			return m.lineForm(false)
		}
		return nil
	}
	if key == "J" || key == "K" {
		return m.reorder(key)
	}
	if key == "x" && m.count() > 0 {
		m.confirm = &confirmation{"Delete selected entry?", func() tea.Cmd {
			data := clone(m.data)
			if m.screen == lists {
				data.Lists = append(data.Lists[:m.cursor], data.Lists[m.cursor+1:]...)
			} else {
				e := data.Lists[m.list].Entries
				data.Lists[m.list].Entries = append(e[:m.cursor], e[m.cursor+1:]...)
			}
			m.data = data
			m.cursor = min(m.cursor, max(0, m.count()-1))
			m.dirty = true
			return m.persist()
		}}
		return nil
	}
	if key != "a" && key != "e" {
		return nil
	}
	if key == "e" && m.count() == 0 {
		return nil
	}
	if m.screen == lists {
		value := ""
		if key == "e" {
			value = m.data.Lists[m.cursor].Name
		}
		return m.openForm("List", []string{"Name"}, []string{value}, func(v []string) (library.Library, error) {
			data := clone(m.data)
			if key == "a" {
				data.Lists = append(data.Lists, library.List{ID: id(), Name: v[0], Settings: library.Defaults()})
			} else {
				data.Lists[m.cursor].Name = v[0]
			}
			return data, nil
		})
	}
	if key == "e" && m.data.Lists[m.list].Entries[m.cursor].Kind != "item" {
		m.notice = "Exchange editing arrives in its dependent slice"
		return nil
	}
	value, spoken := "", ""
	if key == "e" {
		line := m.data.Lists[m.list].Entries[m.cursor].Lines[0]
		value, spoken = line.CW, line.Spoken
	}
	return m.openForm("Practice item · custom speech survives CW edits; review it", []string{"CW text", "Spoken override (clear to use generated speech)"}, []string{value, spoken}, func(v []string) (library.Library, error) {
		data := clone(m.data)
		text := strings.ToUpper(strings.TrimSpace(v[0]))
		e := library.Entry{ID: id(), Kind: "item", Title: text, Lines: []library.Line{{ID: id(), Station: "learner", CW: text, Spoken: v[1]}}}
		if key == "a" {
			data.Lists[m.list].Entries = append(data.Lists[m.list].Entries, e)
		} else {
			old := &data.Lists[m.list].Entries[m.cursor]
			old.Title = text
			old.Lines[0].CW = text
			old.Lines[0].Spoken = v[1]
		}
		return data, nil
	})
}
func (m *Model) reorder(key string) tea.Cmd {
	target := m.cursor + 1
	if key == "K" {
		target = m.cursor - 1
	}
	if target < 0 || target >= m.count() {
		return nil
	}
	data := clone(m.data)
	if m.screen == lists {
		data.Lists[m.cursor], data.Lists[target] = data.Lists[target], data.Lists[m.cursor]
	} else {
		e := data.Lists[m.list].Entries
		e[m.cursor], e[target] = e[target], e[m.cursor]
	}
	m.data = data
	m.cursor = target
	m.dirty = true
	return m.persist()
}

func (m *Model) pickerKey(key string) tea.Cmd {
	f := m.form
	switch key {
	case "esc":
		f.picker = false
	case "j", "down":
		f.pick = (f.pick + 1) % len(morse.Prosigns)
	case "k", "up":
		f.pick = (f.pick + len(morse.Prosigns) - 1) % len(morse.Prosigns)
	case "enter":
		input := &f.fields[f.focus]
		text := []rune(input.Value())
		pos := input.Position()
		token := morse.Prosigns[f.pick].Token
		input.SetValue(string(text[:pos]) + token + string(text[pos:]))
		input.SetCursor(pos + len(token))
		f.picker = false
	}
	return nil
}
func (m *Model) pickerView() string {
	var b strings.Builder
	b.WriteString("Prosigns · sent continuously\n\n")
	for i, p := range morse.Prosigns {
		mark := " "
		if i == m.form.pick {
			mark = ">"
		}
		fmt.Fprintf(&b, "%s %s — %s\n", mark, p.Token, p.Meaning)
	}
	b.WriteString("\nEnter inserts · j/k moves · Esc cancels")
	return b.String()
}
