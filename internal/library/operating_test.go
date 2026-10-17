package library

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPhoneticStarterUsesStandardWordsAndIndividualLetterCW(t *testing.T) {
	want := strings.Fields("Alfa Bravo Charlie Delta Echo Foxtrot Golf Hotel India Juliett Kilo Lima Mike November Oscar Papa Quebec Romeo Sierra Tango Uniform Victor Whiskey X-ray Yankee Zulu")
	list := PhoneticStarter()
	if len(list.Entries) != 26 || list.Settings.Style != "sandwich" {
		t.Fatal("missing phonetic lesson")
	}
	for i, e := range list.Entries {
		line := e.Lines[0]
		if line.CW != string(rune('A'+i)) || line.Spoken != want[i] {
			t.Fatalf("letter %d: %+v", i, line)
		}
	}
}

func TestPhoneticEditsAndRemovalSurviveReload(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "library.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	data, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	index := len(data.Lists) - 1
	data.Lists[index].Entries[0].Lines[0].Spoken = "custom Alfa"
	data.Lists[index].Settings.Display = "hidden"
	data.Lists[index].Entries = data.Lists[index].Entries[:25]
	if err := store.Save(data); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Fatal("starter edits were reseeded")
	}
}
