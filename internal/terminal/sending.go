package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/midi"
	"cw-looper/internal/session"
	"fmt"
)

func (m *Model) startSending() tea.Cmd {
	if !m.midi.connected {
		m.notice = "Connect Vail with m first, then select an item or line and press t"
		return nil
	}
	if m.screen == entries && m.count() > 0 {
		m.entry = m.cursor
		if m.data.Lists[m.list].Entries[m.entry].Kind == "exchange" {
			m.screen = lines
			m.cursor = 0
			m.notice = "Select a line, then t for send-back practice"
			return nil
		}
	} else if m.screen != lines || m.count() == 0 {
		return nil
	}
	m.feedbackOpen = false
	m.preparePractice()
	m.session.SendBack = true
	return m.dispatch(session.Input{Kind: session.Start})
}
func (m *Model) sendingKey(key string) (tea.Cmd, bool) {
	if !m.session.SendBack {
		return nil, false
	}
	switch key {
	case "f":
		if m.session.Phase == session.Captured || m.session.Phase == session.Replay {
			m.feedbackOpen = true
			m.feedbackScroll = 0
		}
		return nil, true
	case "-":
		m.sendWPM = max(5, m.sendWPM-1)
		return nil, true
	case "+", "=":
		m.sendWPM = min(60, m.sendWPM+1)
		return nil, true
	case "enter":
		if m.session.Phase == session.Capturing {
			return m.dispatch(session.Input{Kind: session.FinishCapture, At: midi.HostTime()}), true
		}
		if m.session.Phase == session.Ready || m.session.Phase == session.Captured {
			if !m.midi.connected {
				m.notice = "MIDI disconnected; Esc then m to reconnect"
				return nil, true
			}
			if len(m.midi.down) > 0 {
				m.notice = "Release both pads before starting capture"
				return nil, true
			}
			m.notice = ""
			return m.dispatch(session.Input{Kind: session.BeginCapture, At: midi.HostTime()}), true
		}
	case "y":
		return m.dispatch(session.Input{Kind: session.ReplayAttempt}), true
	case "v":
		return m.dispatch(session.Input{Kind: session.Reference}), true
	case "n":
		return m.dispatch(session.Input{Kind: session.Start}), true
	}
	return nil, false
}
func (m *Model) sendingView() string {
	if m.feedbackOpen {
		return m.feedbackView()
	}
	s := m.session
	title := "Send-back practice"
	text := m.currentLines()[s.Line].CW
	if s.Settings.Display == "hidden" && !s.Revealed {
		text = "Answer hidden · R reveal"
	}
	status := string(s.Phase)
	instruction := "Reference plays once; capture starts only when you press Enter."
	switch s.Phase {
	case session.Ready:
		instruction = "Release pads, then Enter to begin sending."
	case session.Capturing:
		instruction = "Send on the pads · Enter finishes · two-minute limit"
	case session.Captured:
		instruction = "y replay attempt · v reference · Enter new attempt"
	case session.Paused:
		instruction = "Space resumes; paused capture restarts at Ready."
	case session.Failed:
		status = "Playback failed · Enter retries; Esc returns"
	}
	if s.Attempt.Err != nil {
		status += "\n" + s.Attempt.Err.Error()
	}
	if s.Phase == session.Capturing && s.FinishAt > 0 {
		instruction = "Finishing: draining timestamped MIDI input…"
	}
	connection := m.midi.name
	if !m.midi.connected {
		connection = "Disconnected · Esc then m to reconnect"
	}
	return fmt.Sprintf("%s · Line %d of %d\nMIDI: %s\n\n%s\n\n%s\nCaptured marks: %d\n%s\n%s\nAdapter generates element lengths; gaps reflect your spacing.\nn retry reference · Space pause · R reveal · s speak\nh/l previous/next · [/] item · p settings\nEsc library · ? help · q quit", title, s.Line+1, len(m.currentLines()), clip(connection, 65), wrap(text, min(m.width-6, 100)), status, len(s.Attempt.Marks), m.sendingFeedback(), instruction)
}
func (m *Model) midiEvent(event midi.Event) tea.Cmd {
	return m.dispatch(session.Input{Kind: session.KeyEvent, Event: event})
}
func (m *Model) midiLost(err error) tea.Cmd {
	m.notice = err.Error()
	return m.dispatch(session.Input{Kind: session.InputLost, Err: err})
}

func (m *Model) sendingFeedback() string {
	s := m.session
	heading := fmt.Sprintf("Sending %.0f WPM · -/+ adjust · f full feedback", m.sendWPM)
	if s.Phase != session.Captured && s.Phase != session.Replay {
		return heading
	}
	f := s.Attempt.Feedback(m.sendWPM)
	if f.Text == "" {
		return heading + "\n" + f.Status
	}
	comparison := "R reveals target and comparison"
	if s.Settings.Display != "hidden" || s.Revealed {
		comparison = f.Compare(m.currentLines()[s.Line].CW)
	}
	return heading + "\n" + f.Status + ": " + clip(f.Text, max(20, m.width-30)) + "\n" + wrap(f.Timing, max(20, m.width-4)) + "\n" + comparison
}
