package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"path/filepath"
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
