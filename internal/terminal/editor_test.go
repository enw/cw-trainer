package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"path/filepath"
	"strings"
	"testing"
)

func TestExchangeEditorPersistsRolesSpeechAndOrder(t *testing.T) {
	store, err := library.Open(filepath.Join(t.TempDir(), "library.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	m := New(library.Fixture(), output.NewRunner(nil), store.Save)
	key := func(key rune) {
		_, cmd := m.Update(tea.KeyPressMsg{Code: key, Text: string(key)})
		if m.saving && cmd != nil {
			m.Update(cmd())
		}
	}
	save := func(values ...string) {
		t.Helper()
		for i, v := range values {
			m.form.fields[i].SetValue(v)
		}
		_, cmd := m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
		if cmd == nil {
			t.Fatal(m.notice)
		}
		m.Update(cmd())
		if m.dirty {
			t.Fatal(m.notice)
		}
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	key('A')
	save("A full contact", "CQ DE KE2ISH K", "learner", "")
	key('j')
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	key('a')
	save("KE2ISH DE W1AW <AR>", "other", "Your call from W one A W")
	key('j')
	key('K') // Newly added second line moves to first.
	data, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	e := data.Lists[0].Entries[1]
	if len(e.Lines) != 2 || e.Lines[0].Station != "other" || e.Lines[0].Spoken != "Your call from W one A W" {
		t.Fatalf("lost line edit: %+v", e)
	}
	if e.Lines[1].CW != "CQ DE KE2ISH K" {
		t.Fatal("line order lost")
	}
}

func TestLongPhraseAndPasteUseFullCWWithShortTitle(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	phrase := strings.Repeat("CQ DE KE2ISH <AR> ", 20)
	m.Update(tea.PasteMsg{Content: phrase})
	_, cmd := m.Update(tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal(m.notice)
	}
	entry := m.data.Lists[0].Entries[1]
	if entry.Lines[0].CW != strings.TrimSpace(phrase) || len(entry.Title) > 100 {
		t.Fatal("long phrase truncated or rejected")
	}
}

func TestSavingKeepsPracticeControlsResponsive(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.screen = practice
	m.session.Select(2, 0)
	m.session.Phase = session.Paused
	m.saving = true
	m.Update(tea.KeyPressMsg{Code: 'l', Text: "l"})
	if m.session.Line != 1 || m.session.Phase != session.Paused {
		t.Fatal("saving blocked paused navigation")
	}
	m.Update(tea.KeyPressMsg{Code: 'R', Text: "R"})
	if !m.session.Revealed {
		t.Fatal("saving blocked reveal")
	}
}

func TestNormalQuitWaitsForPendingSaveAfterStopping(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	m.saving = true
	m.dirty = true
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if cmd != nil || !m.exitAfterSave {
		t.Fatal("quit abandoned pending save")
	}
	_, cmd = m.Update(savedMsg{})
	if cmd == nil || m.Unsaved() {
		t.Fatal("successful save did not release quit")
	}
}
