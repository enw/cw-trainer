package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"reflect"
	"strings"
	"testing"
)

func TestOperatingReferenceScrollImportAndHiddenPractice(t *testing.T) {
	m := New(library.Fixture(), output.NewRunner(nil))
	before := clone(m.data)
	m.reference = &operatingReference{}
	for i := 0; i < 30; i++ {
		m.referenceKey("j")
	}
	if !strings.Contains(m.referenceView(), "Zulu") {
		t.Fatal("cannot reach alphabet end")
	}
	m.Update(tea.KeyPressMsg{Code: 'i', Text: "i"})
	if m.confirm == nil {
		t.Fatal("import not explicit")
	}
	m.confirm.yes()
	if len(m.data.Lists) != 2 || !reflect.DeepEqual(before.Lists[0], m.data.Lists[0]) {
		t.Fatal("import altered existing material")
	}
	added := m.data.Lists[1]
	if added.ID == library.PhoneticStarter().ID || added.Entries[0].ID == library.PhoneticStarter().Entries[0].ID {
		t.Fatal("import reused starter identity")
	}
	m.referenceKey("esc")
	m.list = 1
	m.entry = 0
	m.screen = practice
	m.session.Select(1, 0)
	m.session.Settings = added.Settings
	m.session.Settings.Display = "hidden"
	m.session.Phase = session.CW
	if strings.Contains(m.View().Content, "Alfa") {
		t.Fatal("phonetic answer leaked")
	}
	m.reference = &operatingReference{}
	m.referenceKey("esc")
	if m.screen != practice || m.session.Revealed {
		t.Fatal("reference changed practice state")
	}
}
