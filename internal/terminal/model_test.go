package terminal_test

import (
	tea "charm.land/bubbletea/v2"
	"context"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/terminal"
	"strings"
	"testing"
)

func TestLibraryNavigationDoesNotStartPlayback(t *testing.T) {
	m := terminal.New(library.Fixture(), output.NewRunner(nil))
	var model tea.Model = m
	for i := 0; i < 2; i++ {
		next, cmd := model.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		if cmd != nil {
			t.Fatal("browsing starts effects")
		}
		model = next
	}
	next, cmd := model.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if cmd != nil {
		t.Fatal("selection starts effects")
	}
	model = next
	view := model.View().Content
	if !strings.Contains(view, "> You") || !strings.Contains(view, "KE2ISH") {
		t.Fatal(view)
	}
}

func TestPausedTopLevelNavigationStaysSilent(t *testing.T) {
	data := library.Fixture()
	second := data.Lists[0].Entries[0]
	second.ID = "second"
	data.Lists[0].Entries = append(data.Lists[0].Entries, second)
	// Use controlled output so only keyboard/session behavior is under test.
	backend := &silentBackend{}
	runner := output.NewRunner(backend)
	defer runner.Close()
	var model tea.Model = terminal.New(data, runner)
	key := func(code rune, text string) tea.Cmd {
		next, cmd := model.Update(tea.KeyPressMsg{Code: code, Text: text})
		model = next
		return cmd
	}
	key(tea.KeyEnter, "")
	key(tea.KeyEnter, "")
	key(tea.KeyEnter, "")
	stop := key(' ', " ")
	// Execute the cancellation callback and deliver its acknowledgment.
	if stop != nil {
		model, _ = model.Update(stop())
	}
	key(']', "]")
	if !strings.Contains(model.View().Content, "Paused") {
		t.Fatal("next item resumed paused session: " + model.View().Content)
	}
}

type silentBackend struct{}

func (*silentBackend) CW(ctx context.Context, text string, hz float64) error {
	<-ctx.Done()
	return ctx.Err()
}
func (*silentBackend) Speech(ctx context.Context, text string) error { <-ctx.Done(); return ctx.Err() }

func TestResizeStopsPendingExchangeSwitch(t *testing.T) {
	model, runner := twoExchanges()
	defer runner.Close()
	for i := 0; i < 3; i++ {
		press(model, "enter")
	}
	stopped := press(model, "]")
	model.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	model.Update(stopped())
	model.Update(tea.WindowSizeMsg{Width: 100, Height: 28})
	if strings.Contains(model.View().Content, "CW · Completed") {
		t.Fatal("resize restarted pending playback")
	}
}
func TestRapidExchangeNavigationKeepsPausedSpeechIntent(t *testing.T) {
	model, runner := twoExchanges()
	defer runner.Close()
	for i := 0; i < 3; i++ {
		press(model, "enter")
	}
	pause := press(model, "space")
	model.Update(pause())
	press(model, "s")
	stopped := press(model, "]")
	press(model, "[")
	model.Update(stopped())
	if !strings.Contains(model.View().Content, "Paused") {
		t.Fatal("rapid navigation resumed a paused session")
	}
}
func TestSpeakDuringExchangeSwitchUsesNewTarget(t *testing.T) {
	model, runner := twoExchanges()
	defer runner.Close()
	for i := 0; i < 3; i++ {
		press(model, "enter")
	}
	stopped := press(model, "]")
	press(model, "s")
	model.Update(stopped())
	view := model.View().Content
	if !strings.Contains(view, "Second exchange") || !strings.Contains(view, "Speaking") {
		t.Fatal("new target stranded: " + view)
	}
}
func twoExchanges() (*terminal.Model, *output.Runner) {
	data := library.Fixture()
	second := data.Lists[0].Entries[0]
	second.ID = "second"
	second.Title = "Second exchange"
	data.Lists[0].Entries = append(data.Lists[0].Entries, second)
	runner := output.NewRunner(&silentBackend{})
	return terminal.New(data, runner), runner
}
func press(model *terminal.Model, key string) tea.Cmd {
	var msg tea.KeyPressMsg
	switch key {
	case "enter":
		msg.Code = tea.KeyEnter
	case "space":
		msg.Code = ' '
	default:
		msg.Code = []rune(key)[0]
		msg.Text = key
	}
	_, cmd := model.Update(msg)
	return cmd
}
