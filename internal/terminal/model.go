// Package terminal maps keyboard actions to the practice session.
package terminal

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
)

type screen int

const (
	lists screen = iota
	entries
	lines
	practice
)

type Model struct {
	pendingPaused       bool
	data                library.Library
	runner              *output.Runner
	session             *session.Session
	screen              screen
	list, entry, cursor int
	width, height       int
	pendingEntry        *int
	help                bool
}

func New(data library.Library, runner *output.Runner) *Model {
	return &Model{data: data, runner: runner, session: session.New(0), width: 80, height: 24}
}
func (m *Model) Init() tea.Cmd { return nil }
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case session.Input:
		return m, m.dispatch(msg)
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.small() && m.screen == practice {
			return m, m.dispatch(session.Input{Kind: session.Back})
		}
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			m.pendingEntry = nil
			return m, m.dispatch(session.Input{Kind: session.Quit})
		}
		if m.small() {
			return m, nil
		}
		if key == "?" {
			m.help = !m.help
			return m, nil
		}
		if m.help {
			if key == "esc" {
				m.help = false
			}
			return m, nil
		}
		if m.screen == practice {
			return m, m.practiceKey(key)
		}
		m.browseKey(key)
		if m.screen == practice {
			return m, m.dispatch(session.Input{Kind: session.Start})
		}
	}
	return m, nil
}
func (m *Model) small() bool                  { return m.width < 80 || m.height < 24 }
func (m *Model) currentLines() []library.Line { return m.data.Lists[m.list].Entries[m.entry].Lines }
func (m *Model) dispatch(in session.Input) tea.Cmd {
	effects := m.session.Update(in, time.Now())
	var cmds []tea.Cmd
	exiting := false
	for _, effect := range effects {
		if effect.Kind == session.Exit {
			cmds = append(cmds, tea.Quit)
			exiting = true
			continue
		}
		var line library.Line
		if effect.Kind == session.PlayCW || effect.Kind == session.PlaySpeech {
			line = m.currentLines()[m.session.Line]
		}
		events := m.runner.Execute(effect, line)
		cmds = append(cmds, func() tea.Msg { return <-events })
	}
	if m.session.Phase == session.Idle && m.screen == practice && !exiting {
		if m.pendingEntry != nil {
			m.entry = *m.pendingEntry
			m.pendingEntry = nil
			m.cursor = 0
			m.session.Select(len(m.currentLines()), 0)
			kind := session.Start
			if m.pendingPaused {
				kind = session.StartPaused
			}
			m.pendingPaused = false
			cmds = append(cmds, m.dispatch(session.Input{Kind: kind}))
		} else {
			m.screen = lines
			m.cursor = m.session.Line
		}
	}
	return tea.Batch(cmds...)
}
func (m *Model) practiceKey(key string) tea.Cmd {
	action := session.Input{}
	switch key {
	case "space":
		action.Kind = session.TogglePause
	case "s", "S":
		action.Kind = session.Speak
	case "left", "h":
		action.Kind = session.Previous
	case "right", "l":
		action.Kind = session.Next
	case "esc":
		m.pendingEntry = nil
		action.Kind = session.Back
	case "enter":
		if m.session.Phase != session.Failed {
			return nil
		}
		action.Kind = session.Start
	case "[", "]":
		step := 1
		if key == "[" {
			step = -1
		}
		current := m.entry
		if m.pendingEntry != nil {
			current = *m.pendingEntry
		}
		target := min(max(current+step, 0), len(m.data.Lists[m.list].Entries)-1)
		if target == current {
			return nil
		}
		m.pendingPaused = m.session.WantsPause()
		m.pendingEntry = &target
		action.Kind = session.Back
	default:
		return nil
	}
	return m.dispatch(action)
}
func (m *Model) browseKey(key string) {
	count := m.count()
	switch key {
	case "up", "k":
		m.cursor = max(0, m.cursor-1)
	case "down", "j":
		m.cursor = min(max(0, count-1), m.cursor+1)
	case "esc", "left", "h":
		switch m.screen {
		case entries:
			m.screen = lists
			m.cursor = m.list
		case lines:
			m.screen = entries
			m.cursor = m.entry
		}
	case "enter", "right", "l":
		if count == 0 {
			return
		}
		switch m.screen {
		case lists:
			m.list = m.cursor
			m.screen = entries
			m.cursor = 0
		case entries:
			m.entry = m.cursor
			m.screen = lines
			m.cursor = 0
		case lines:
			m.session.Select(count, m.cursor)
			m.screen = practice
		}
	}
}
func (m *Model) count() int {
	switch m.screen {
	case lists:
		return len(m.data.Lists)
	case entries:
		return len(m.data.Lists[m.list].Entries)
	case lines:
		return len(m.currentLines())
	}
	return 0
}
func (m *Model) View() tea.View {
	var content string
	switch {
	case m.small():
		content = "cw-looper\n\nResize the terminal to at least 80 columns × 24 rows.\nPlayback stopped. Selection is preserved.\n\nq quit"
	case m.help:
		content = "cw-looper · Help\n\nBrowse: arrows or j/k select; Enter/l opens; Esc/h goes back.\nPractice: Space pauses; h/l or arrows change lines; [/] change items.\ns speaks the current phrase; Esc stops and returns; q quits.\n\nPause restarts interrupted CW/speech; gaps resume their remainder.\nNavigation stays silent while paused. Lines repeat until you advance.\n\nThis skeleton has fixed 20/15 WPM, 440/600 Hz and 2-second gaps.\nEditing and other modes arrive in later slices.\n\n? or Esc closes help · q quit"
	case m.screen == practice:
		content = m.practiceView()
	default:
		content = m.browseView()
	}
	view := tea.NewView(content)
	view.AltScreen = true
	return view
}
func (m *Model) browseView() string {
	title := "Lists"
	var rows []string
	switch m.screen {
	case lists:
		for _, l := range m.data.Lists {
			rows = append(rows, l.Name)
		}
	case entries:
		title = m.data.Lists[m.list].Name
		for _, e := range m.data.Lists[m.list].Entries {
			rows = append(rows, e.Title+" (exchange)")
		}
	case lines:
		title = m.data.Lists[m.list].Entries[m.entry].Title
		for _, l := range m.currentLines() {
			role, _ := station(l)
			rows = append(rows, role+" · "+l.CW)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "cw-looper · Saved\n\n%s\n\n", title)
	if len(rows) == 0 {
		b.WriteString("No saved material. Editing is not part of this skeleton.\n")
	}
	start := max(0, m.cursor-(m.height-10)+1)
	for i := start; i < min(len(rows), start+m.height-10); i++ {
		marker := " "
		if i == m.cursor {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s %s\n", marker, clip(rows[i], m.width-4))
	}
	b.WriteString("\nj/k or arrows select · Enter/l open · Esc/h back · ? help · q quit")
	return b.String()
}
func (m *Model) practiceView() string {
	line := m.currentLines()[m.session.Line]
	role, hz := station(line)
	text := wrap(line.CW, min(m.width-6, 100))
	status := string(m.session.Phase)
	if m.session.Phase == session.Paused {
		status += " · Space restarts the interrupted portion"
	}
	if m.session.Phase == session.Failed {
		status = wrap(m.session.Err.Error(), m.width-6) + "\nEnter retries · restart app if indicated"
	}
	boundary := ""
	if m.session.Line == len(m.currentLines())-1 {
		boundary = "Last line · ] next item or Esc library\n"
	}
	return fmt.Sprintf("cw-looper · %s · Line %d of %d\nStation: %s · %d Hz\n\n%s\n\n%s · Completed: %d\n20 character / 15 overall WPM · Pause 2 s\nWhole line · CW only · Repeat current\n\n%sSpace pause · h/l or arrows previous/next line\n[/] previous/next item · s speak · Esc library · ? help · q quit",
		clip(m.data.Lists[m.list].Entries[m.entry].Title, 40), m.session.Line+1, len(m.currentLines()), role, hz, text, status, m.session.Count, boundary)
}
func station(line library.Line) (string, int) {
	if line.Station == "other" {
		return "Other", 600
	}
	return "You", 440
}
func clip(s string, width int) string {
	r := []rune(s)
	if len(r) > width {
		return string(r[:max(0, width-1)]) + "…"
	}
	return s
}
func wrap(s string, width int) string {
	var lines []string
	r := []rune(s)
	for len(r) > width {
		lines = append(lines, string(r[:width]))
		r = r[width:]
	}
	return strings.Join(append(lines, string(r)), "\n")
}
