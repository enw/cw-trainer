package library

import (
	"cw-looper/internal/morse"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRSTReportsConnectRatingsToSpeechAndOrderedStations(t *testing.T) {
	for _, bad := range []string{"59", "699", "500", "5N9", "5999"} {
		if _, err := ReportSpeech(bad); err == nil {
			t.Fatalf("invalid report accepted: %s", bad)
		}
	}
	speech, err := ReportSpeech("339")
	if err != nil || !strings.Contains(speech, "Readability 3") || !strings.Contains(speech, "Strength 3") || !strings.Contains(speech, "Tone 9") {
		t.Fatal(speech, err)
	}
	list := RSTStarter()
	for _, report := range []string{"599", "579", "339"} {
		found := false
		for _, e := range list.Entries {
			if e.Kind == "item" && e.Lines[0].CW == report {
				found = true
			}
		}
		if !found {
			t.Fatal("missing report", report)
		}
	}
	exchange := list.Entries[len(list.Entries)-1]
	if exchange.Kind != "exchange" || len(exchange.Lines) != 2 || exchange.Lines[0].Station != "learner" || exchange.Lines[1].Station != "other" {
		t.Fatal("station order")
	}
	for _, line := range exchange.Lines {
		if _, err := morse.Compile(line.CW, 20, 15, false); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(line.CW, "KE2ISH") || !strings.Contains(line.CW, "W1AW") || line.Spoken == "" {
			t.Fatal(line)
		}
	}
}

func TestRSTScaleRangesAndEditedReportsSurviveReload(t *testing.T) {
	scales := RSTScales()
	if len(scales) != 3 || scales[0].Name != "Readability" || len(scales[0].Meanings) != 5 || scales[1].Name != "Strength" || len(scales[1].Meanings) != 9 || scales[2].Name != "Tone" || len(scales[2].Meanings) != 9 {
		t.Fatal("RST ranges/order")
	}
	store, err := Open(filepath.Join(t.TempDir(), "library.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	data, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	list := &data.Lists[5]
	list.Entries[0].Lines[0].Spoken = "my report prompt"
	list.Settings.LineAdvancement = "automatic"
	list.Settings.Display = "hidden"
	list.Entries = append(list.Entries[:1], list.Entries[2:]...)
	if err := store.Save(data); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Fatal("RST customization replaced")
	}
}
