package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/midi"
	"cw-looper/internal/session"
	"fmt"
	"strings"
	"time"
)

type midiState struct {
	manager               *midi.Manager
	open, busy, connected bool
	devices               []midi.Device
	cursor                int
	generation            uint64
	name, status          string
	recent                []string
	down                  map[byte]time.Duration
	received              int
}
type midiDevicesMsg struct {
	devices []midi.Device
	err     error
}
type midiConnectedMsg struct {
	name string
	err  error
}
type midiPollMsg struct {
	finishAt   time.Duration
	generation uint64
	events     []midi.Event
	err        error
}
type midiTickMsg struct{ generation uint64 }

func (m *Model) CloseMIDI() error { return m.midi.manager.Close() }
func (m *Model) openMIDI() tea.Cmd {
	m.midi.open = true
	m.midi.busy = true
	return func() tea.Msg { devices, err := m.midi.manager.Devices(); return midiDevicesMsg{devices, err} }
}
func midiTick(generation uint64) tea.Cmd {
	return tea.Tick(20*time.Millisecond, func(time.Time) tea.Msg { return midiTickMsg{generation} })
}
func (m *Model) midiMessage(msg tea.Msg) tea.Cmd {
	state := &m.midi
	switch msg := msg.(type) {
	case midiDevicesMsg:
		state.busy = false
		state.devices = msg.devices
		state.cursor = 0
		if msg.err != nil {
			state.status = msg.err.Error()
		} else if len(msg.devices) == 0 {
			state.status = "No MIDI sources. Connect adapter, then f to refresh."
		}
	case midiConnectedMsg:
		state.busy = false
		state.connected = msg.err == nil
		state.generation++
		if msg.err != nil {
			state.status = msg.err.Error()
			return nil
		}
		state.name = msg.name
		state.status = "MIDI enabled; waiting for keying events"
		state.recent = nil
		state.received = 0
		state.down = make(map[byte]time.Duration)
		return midiTick(state.generation)
	case midiTickMsg:
		if !state.connected || msg.generation != state.generation {
			return nil
		}
		finishAt := m.session.FinishAt
		if midi.HostTime() < finishAt+100*time.Millisecond {
			finishAt = 0
		}
		return func() tea.Msg {
			events, err := state.manager.Poll()
			return midiPollMsg{generation: msg.generation, events: events, err: err, finishAt: finishAt}
		}
	case midiPollMsg:
		if !state.connected || msg.generation != state.generation {
			return nil
		}
		if msg.err != nil {
			state.connected = false
			state.generation++
			state.down = nil
			state.status = msg.err.Error()
			return m.midiLost(msg.err)
		}
		commands := []tea.Cmd{}
		for _, event := range msg.events {
			state.received++
			label := "up"
			duration := ""
			if event.Down {
				label = "down"
				state.down[event.Note] = event.At
			} else if start, ok := state.down[event.Note]; ok {
				duration = fmt.Sprintf(" · %.1f ms", float64(event.At-start)/float64(time.Millisecond))
				delete(state.down, event.Note)
			}
			state.recent = append(state.recent, fmt.Sprintf("%.3f s · note %d %s%s", event.At.Seconds(), event.Note, label, duration))
			if len(state.recent) > 8 {
				state.recent = state.recent[len(state.recent)-8:]
			}
			commands = append(commands, m.midiEvent(event))
		}
		if state.received > 0 {
			state.status = fmt.Sprintf("Receiving · %d events · %d held notes", state.received, len(state.down))
		}
		if msg.finishAt > 0 {
			commands = append(commands, m.dispatch(session.Input{Kind: session.CaptureDrained, At: msg.finishAt}))
		}
		commands = append(commands, m.dispatch(session.Input{Kind: session.CaptureClock, At: midi.HostTime()}), midiTick(state.generation))
		return tea.Batch(commands...)
	}
	return nil
}
func (m *Model) midiKey(key string) tea.Cmd {
	state := &m.midi
	if key == "esc" {
		state.open = false
		return nil
	}
	if state.busy {
		return nil
	}
	switch key {
	case "f":
		return m.openMIDI()
	case "j", "down":
		state.cursor = min(max(0, len(state.devices)-1), state.cursor+1)
	case "k", "up":
		state.cursor = max(0, state.cursor-1)
	case "d":
		state.connected = false
		state.generation++
		state.down = nil
		state.status = "Disconnected; keyboard output restored"
		if err := state.manager.Disconnect(); err != nil {
			state.status = err.Error()
		}
	case "enter":
		if len(state.devices) == 0 {
			return nil
		}
		device := state.devices[state.cursor]
		state.busy = true
		state.connected = false
		state.generation++
		return func() tea.Msg { return midiConnectedMsg{device.Name, state.manager.Connect(device.ID)} }
	}
	return nil
}
func (m *Model) midiView() string {
	state := &m.midi
	var b strings.Builder
	b.WriteString("Vail USB MIDI · device connection and input diagnostic\n\n")
	start := max(0, state.cursor-3)
	for i := start; i < min(len(state.devices), start+4); i++ {
		device := state.devices[i]
		mark := " "
		if i == state.cursor {
			mark = ">"
		}
		fmt.Fprintf(&b, "%s %s\n", mark, clip(device.Name, 65))

	}
	fmt.Fprintf(&b, "\n%s\n", wrap(state.status, 74))
	for _, line := range state.recent {
		b.WriteString(line + "\n")
	}
	if state.busy {
		b.WriteString("Working…\n")
	}
	b.WriteString("\nEnter enables MIDI on selected device · f refresh · d disconnect\nIambic keyer, speed and sidetone unchanged; quit restores keyboard mode.\nEsc returns with MIDI connected · j/k select · q quit")
	return b.String()
}
