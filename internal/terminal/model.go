// Package terminal maps keyboard actions to the practice session.
package terminal

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/midi"
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
	midi                     midiState
	data                     library.Library
	runner                   *output.Runner
	session                  *session.Session
	screen                   screen
	list, entry, cursor      int
	width, height            int
	help                     bool
	form                     *form
	confirm                  *confirmation
	save                     func(library.Library) error
	saving, dirty            bool
	exitAfterSave, forceExit bool
	notice                   string
}

func New(data library.Library, runner *output.Runner, save ...func(library.Library) error) *Model {
	m := &Model{data: data, runner: runner, session: session.New(0), width: 80, height: 24}
	m.midi.manager = midi.NewManager(midi.Mac{})
	if len(save) > 0 {
		m.save = save[0]
	}
	return m
}
func (m *Model) Init() tea.Cmd { return nil }
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case midiDevicesMsg, midiConnectedMsg, midiPollMsg, midiTickMsg:
		return m, m.midiMessage(msg)
	case progressMsg:
		if msg.id != m.session.CurrentID() || m.session.Phase != session.CW {
			return m, nil
		}
		id, sample := m.runner.Progress()
		m.session.Update(session.Input{Kind: session.Progress, ID: id, Sample: sample}, time.Now())
		return m, progressTick(msg.id)
	case savedMsg:
		m.saving = false
		if msg.err != nil {
			m.notice = "Not saved: " + msg.err.Error() + " · Ctrl+S retry"
		} else {
			m.dirty = false
			m.notice = "Saved"
		}
		if m.exitAfterSave {
			m.exitAfterSave = false
			return m, m.finishExit()
		}
		return m, nil
	case session.Input:
		return m, m.dispatch(msg)
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.small() && m.screen == practice {
			return m, m.dispatch(session.Input{Kind: session.Back})
		}
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" {
			m.forceExit = true
			if m.exitAfterSave {
				return m, tea.Quit
			}
			return m, m.dispatch(session.Input{Kind: session.Quit})
		}
		if m.saving {
			if key == "q" {
				return m, m.dispatch(session.Input{Kind: session.Quit})
			}
			if m.screen == practice {
				return m, m.practiceKey(key)
			}
			return m, nil
		}
		if m.form != nil {
			return m, m.formKey(msg)
		}
		if m.confirm != nil {
			return m, m.confirmKey(key)
		}
		if key == "i" && m.screen == lists {
			return m, m.importStarter()
		}
		if key == "ctrl+s" && m.dirty {
			return m, m.persist()
		}
		if key == "q" && m.dirty {
			m.confirm = &confirmation{"Discard unsaved changes and quit?", func() tea.Cmd { m.forceExit = true; return m.dispatch(session.Input{Kind: session.Quit}) }}
			return m, nil
		}
		if key == "q" || key == "ctrl+c" {
			return m, m.dispatch(session.Input{Kind: session.Quit})
		}
		if m.midi.open {
			return m, m.midiKey(key)
		}
		if key == "m" && m.screen != practice {
			return m, m.openMIDI()
		}
		if m.small() {
			return m, nil
		}
		if key == "p" && m.screen != lists {
			return m, m.settingsForm()
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
		if key == "t" {
			return m, m.startSending()
		}
		if key == "a" || key == "A" || key == "e" || key == "x" || key == "J" || key == "K" {
			return m, m.editKey(key)
		}
		m.browseKey(key)
		if m.screen == practice {
			return m, m.dispatch(session.Input{Kind: session.Start})
		}
	default:
		if m.form != nil && !m.form.picker {
			f := m.form
			var cmd tea.Cmd
			f.fields[f.focus], cmd = f.fields[f.focus].Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
func (m *Model) small() bool                  { return m.width < 80 || m.height < 24 }
func (m *Model) currentLines() []library.Line { return m.data.Lists[m.list].Entries[m.entry].Lines }
func (m *Model) dispatch(in session.Input) tea.Cmd {
	return m.applyEffects(m.session.Update(in, time.Now()))
}
func (m *Model) applyEffects(effects []session.Effect) tea.Cmd {
	if m.screen == practice {
		m.entry = m.session.Entry
	}
	var cmds []tea.Cmd
	exiting := false
	for _, effect := range effects {
		if effect.Kind == session.Exit {
			cmds = append(cmds, m.finishExit())
			exiting = true
			continue
		}
		var line library.Line
		if effect.Kind == session.PlayCW || effect.Kind == session.PlaySpeech || effect.Kind == session.PlayAttempt {
			line = m.currentLines()[m.session.Line]
		}
		events := m.runner.Execute(effect, line)
		if effect.Kind == session.PlayCW {
			cmds = append(cmds, progressTick(effect.ID))
		}
		cmds = append(cmds, func() tea.Msg { return <-events })
	}
	if m.session.Phase == session.Idle && m.screen == practice && !exiting {
		m.screen = lines
		m.cursor = m.session.Line
		if m.data.Lists[m.list].Entries[m.entry].Kind == "item" {
			m.screen = entries
			m.cursor = m.entry
		}
	}
	return tea.Batch(cmds...)
}
func (m *Model) practiceKey(key string) tea.Cmd {
	if cmd, handled := m.sendingKey(key); handled {
		return cmd
	}
	action := session.Input{}
	switch key {
	case "space":
		action.Kind = session.TogglePause
	case "r", "R":
		action.Kind = session.Reveal
	case "s", "S":
		action.Kind = session.Speak
	case "left", "h":
		action.Kind = session.Previous
		if m.data.Lists[m.list].Entries[m.entry].Kind == "item" {
			action.Kind = session.PreviousEntry
		}
	case "right", "l":
		action.Kind = session.Next
		if m.data.Lists[m.list].Entries[m.entry].Kind == "item" {
			action.Kind = session.NextEntry
		}
	case "esc":
		action.Kind = session.Back
	case "enter":
		if m.session.Phase != session.Failed {
			return nil
		}
		action.Kind = session.Start
	case "[":
		action.Kind = session.PreviousEntry
	case "]":
		action.Kind = session.NextEntry
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
			if m.data.Lists[m.list].Entries[m.entry].Kind == "item" {
				m.preparePractice()
			}
		case lines:
			m.preparePractice()
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

func (m *Model) preparePractice() {
	m.notice = ""
	lengths := make([]int, len(m.data.Lists[m.list].Entries))
	exchanges := make([]bool, len(lengths))
	for i, e := range m.data.Lists[m.list].Entries {
		lengths[i] = len(e.Lines)
		exchanges[i] = e.Kind == "exchange"
	}
	m.session.SelectEntries(lengths, m.entry, m.cursor)
	m.session.SetExchanges(exchanges)
	m.session.Settings = m.data.Lists[m.list].Settings
	m.screen = practice
}

type progressMsg struct{ id uint64 }

func progressTick(id uint64) tea.Cmd {
	return tea.Tick(25*time.Millisecond, func(time.Time) tea.Msg { return progressMsg{id} })
}

func (m *Model) Unsaved() bool { return m.dirty || m.saving }
func (m *Model) Saving() bool  { return m.saving }

// Both save completion and output cleanup can finish first. Decide exit only
// here, after output is stopped, so failure never silently discards edits.
func (m *Model) finishExit() tea.Cmd {
	if m.forceExit {
		return tea.Quit
	}
	if m.saving {
		m.exitAfterSave = true
		return nil
	}
	if !m.dirty {
		return tea.Quit
	}
	m.session.CancelQuit()
	return m.dispatch(session.Input{Kind: session.Back})
}
