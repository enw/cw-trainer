package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/morse"
	"cw-looper/internal/session"
	"fmt"
	"strings"
)

func (m *Model) View() tea.View {
	var content string
	switch {
	case m.small():
		content = "cw-looper\n\nResize the terminal to at least 80 columns × 24 rows.\nPlayback stopped. Selection is preserved.\n\nq quit"
	case m.form != nil:
		content = m.formView()
	case m.confirm != nil:
		content = m.confirm.prompt + "\n\ny confirm · Esc/n cancel (default)"
	case m.help:
		content = "cw-looper · Help\n\nBrowse: arrows or j/k select; Enter/l opens; Esc/h goes back.\nPractice: Space pauses; h/l or arrows change lines; [/] change items.\ns speaks the current phrase; Esc stops and returns; q quits.\n\nPause restarts interrupted CW/speech; gaps resume their remainder.\nNavigation stays silent while paused. Lines repeat until you advance.\n\nThis skeleton has fixed 20/15 WPM, 440/600 Hz and 2-second gaps.\nEditing and other modes arrive in later slices.\n\n? or Esc closes help · q quit"
	case m.screen == practice:
		content = m.practiceView()
	default:
		content = m.browseView()
	}
	if m.notice != "" {
		content += "\n\n" + m.notice
	}
	if m.saving {
		content += "\nSaving…"
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
			rows = append(rows, e.Title+" ("+e.Kind+")")
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
		b.WriteString("No saved material. Press a to add.\n")
	}
	start := max(0, m.cursor-(m.height-10)+1)
	for i := start; i < min(len(rows), start+m.height-10); i++ {
		marker := " "
		if i == m.cursor {
			marker = ">"
		}
		fmt.Fprintf(&b, "%s %s\n", marker, clip(rows[i], m.width-4))
	}
	b.WriteString("\nj/k select · Enter open · a item/line · A exchange · e edit · x delete · J/K reorder\nEsc/h back · ? help · q quit")
	return b.String()
}
func (m *Model) practiceView() string {
	line := m.currentLines()[m.session.Line]
	role, _ := station(line)
	config := m.session.ActiveSettings
	hz := config.LearnerHz
	if line.Station == "other" {
		hz = config.OtherHz
	}
	text := wrap(line.CW, min(m.width-6, 100))
	title := m.data.Lists[m.list].Entries[m.entry].Title
	if m.session.Settings.Display == "hidden" && !m.session.Revealed {
		text = "Answer hidden · R reveal"
		title = "Listening practice"
	} else if m.session.Settings.Display == "highlight" && m.session.Phase == session.CW {
		timeline, err := morse.Compile(line.CW, config.CharacterWPM, config.OverallWPM, false)
		if err == nil {
			for _, marker := range timeline.Markers {
				if m.session.Sample >= marker.Start && m.session.Sample < marker.End {
					text = wrap(timeline.Text[:marker.From]+"["+timeline.Text[marker.From:marker.To]+"]"+timeline.Text[marker.To:], min(m.width-6, 100))
					break
				}
			}
		}
	}
	status := string(m.session.Phase)
	if m.session.Settings != config {
		status += " · Settings pending next CW"
	}
	if m.session.Phase == session.Paused {
		status += " · Space restarts the interrupted portion"
	}
	if m.session.Phase == session.Failed && m.session.Settings.Display == "hidden" && !m.session.Revealed {
		status = "Playback failed · Enter retries; Esc for details in library"
	} else if m.session.Phase == session.Failed {
		status = wrap(m.session.Err.Error(), m.width-6) + "\nEnter retries · restart app if indicated"
	}
	boundary := ""
	if m.data.Lists[m.list].Entries[m.entry].Kind == "exchange" {
		boundary = "Automatic lines: ordered progression\n"
		if m.session.Settings.LineAdvancement == "manual" {
			boundary = "Manual lines: press Next\n"
		}
	}
	if m.session.Line == len(m.currentLines())-1 {
		boundary += "Last line · ] next item or Esc library\n"
	}
	return fmt.Sprintf("cw-looper · %s · Line %d of %d\nStation: %s · %.0f Hz\n\n%s\n\n%s · Completed: %d\n%g character / %g overall WPM · Pause %g s\n%s · %s · %s\n\n%sSpace pause · h/l or arrows previous/next line\n[/] previous/next item · s speak · R reveal · p settings · Esc library · ? help · q quit",
		clip(title, 40), m.session.Line+1, len(m.currentLines()), role, hz, text, status, m.session.Count, config.CharacterWPM, config.OverallWPM, m.session.Settings.PauseSeconds, m.session.Settings.Display, m.session.Settings.Style, m.session.Settings.Sequence, boundary)
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
