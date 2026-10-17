package terminal

import (
	"cw-looper/internal/library"
	"cw-looper/internal/morse"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"strings"
	"testing"
	"time"
)

func TestHiddenViewAndProsignProgress(t *testing.T) {
	data := library.Fixture()
	data.Lists[0].Entries[0].Lines[0].CW = "CQ <AR>"
	data.Lists[0].Entries[0].Title = "SECRET ANSWER"
	m := New(data, output.NewRunner(nil))
	m.screen = practice
	m.session.Select(2, 0)
	m.session.Settings.Display = "hidden"
	m.session.Phase = session.CW
	view := m.View().Content
	if strings.Contains(view, "SECRET") || strings.Contains(view, "CQ") || strings.Contains(view, "<AR>") {
		t.Fatal("answer leaked: " + view)
	}
	m.session.Update(session.Input{Kind: session.Reveal}, time.Now())
	if !strings.Contains(m.View().Content, "CQ <AR>") {
		t.Fatal("reveal failed")
	}
	m.session.Settings.Display = "highlight"
	timeline, err := morse.Compile("CQ <AR>", 20, 15, false)
	if err != nil {
		t.Fatal(err)
	}
	marker := timeline.Markers[len(timeline.Markers)-1]
	m.session.Update(session.Input{Kind: session.Progress, ID: m.session.CurrentID(), Sample: marker.Start}, time.Now())
	if !strings.Contains(m.View().Content, "[<AR>]") {
		t.Fatal("prosign was not highlighted as one unit")
	}
}
