package library

import (
	"path/filepath"
	"testing"
)

func TestStarterAndUpgradePreserveUserLibrary(t *testing.T) {
	if err := Starter().Validate(); err != nil {
		t.Fatal(err)
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
	if len(data.Lists) != 4 || len(data.Lists[3].Entries) != 6 {
		t.Fatal("starter categories missing")
	}
	data.Lists[0].Settings.Style = "occasional"
	data.Lists[0].Settings.Display = "hidden"
	data.Lists[0].Settings.Sequence = "shuffle"
	data.Lists[0].Settings.LearnerHz = 500
	data.Lists[1].Entries = nil
	if err := store.Save(data); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Lists[0].Settings != data.Lists[0].Settings || len(got.Lists[1].Entries) != 0 {
		t.Fatal("settings or deliberate empty list replaced")
	}
	// An existing walking-skeleton library is preserved, including its seed version.
	old := Fixture()
	if err := store.Save(old); err != nil {
		t.Fatal(err)
	}
	got, err = store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.SeedVersion != 1 || len(got.Lists) != 1 || got.Lists[0].ID != "starter" {
		t.Fatal("upgrade overwrote old library")
	}
}
