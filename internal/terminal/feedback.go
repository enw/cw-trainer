package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/session"
	"fmt"
	"strings"
)

func (m *Model) feedbackKey(key string) tea.Cmd {
	switch key {
	case "space":
		return m.dispatch(session.Input{Kind: session.TogglePause})
	case "esc", "f":
		m.feedbackOpen = false
	case "j", "down":
		m.feedbackScroll = min(m.feedbackScroll+1, max(0, len(m.feedbackRows())-max(1, m.height-6)))
	case "k", "up":
		m.feedbackScroll = max(0, m.feedbackScroll-1)
	case "R", "r":
		return m.dispatch(session.Input{Kind: session.Reveal})
	case "-", "+", "=":
		cmd, _ := m.sendingKey(key)
		return cmd
	}
	return nil
}

func (m *Model) feedbackRows() []string {
	f := m.session.Attempt.Feedback(m.sendWPM)
	lines := []string{f.Status, "Your sending: " + f.Text, f.Timing}
	if m.session.Settings.Display == "hidden" && !m.session.Revealed {
		lines = append(lines, "Target hidden · R reveals target and comparison")
	} else {
		target := m.currentLines()[m.session.Line].CW
		lines = append(lines, "Target: "+target, f.Compare(target))
	}
	lines = append(lines, "Adapter generates element lengths; this measures output timing, not raw paddle contact.", "Gap ideals: element 1, character 3, word 7 dit units. Other timing is marked uncertain; sending WPM must match the adapter.")
	var rows []string
	for _, line := range lines {
		rows = append(rows, strings.Split(wrap(line, max(20, m.width-4)), "\n")...)
	}
	return rows
}

func (m *Model) feedbackView() string {
	rows := m.feedbackRows()
	height := max(1, m.height-6)
	offset := min(m.feedbackScroll, max(0, len(rows)-height))
	end := min(len(rows), offset+height)
	return fmt.Sprintf("Sending feedback · %.0f WPM · rows %d–%d/%d\n\n%s\n\nj/k scroll · -/+ sending WPM · R reveal\nSpace pause/resume · f/Esc back · q quit", m.sendWPM, offset+1, end, len(rows), strings.Join(rows[offset:end], "\n"))
}
