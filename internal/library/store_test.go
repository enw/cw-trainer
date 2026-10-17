package library_test

import (
	"cw-looper/internal/library"
	"os"
	"path/filepath"
	"testing"
)

func TestSeedRoundTripAndLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "library.json")
	store, err := library.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if data.Lists[0].Entries[0].Lines[1].CW != "W1AW DE KE2ISH K" {
		t.Fatal("wrong fixture")
	}
	if second, err := library.Open(path); err == nil {
		second.Close()
		t.Fatal("concurrent writer admitted")
	}
	data.Lists[0].Entries[0].Lines[1].CW = "KE2ISH"
	if err := store.Save(data); err != nil {
		t.Fatal(err)
	}
	store.Close()
	store, err = library.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	data, err = store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if data.Lists[0].Entries[0].Lines[1].CW != "KE2ISH" {
		t.Fatal("edit overwritten")
	}
	data.Lists = nil
	if err := store.Save(data); err != nil {
		t.Fatal(err)
	}
	data, err = store.Load()
	if err != nil || len(data.Lists) != 0 {
		t.Fatalf("empty library reseeded: %v", err)
	}
}
func TestLoadFailureNeverSeeds(t *testing.T) {
	for _, input := range []string{"broken", `{"schemaVersion":99}`, `{"schemaVersion":1,"unknown":true}`} {
		t.Run(input, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "library.json")
			os.WriteFile(path, []byte(input), 0600)
			store, err := library.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer store.Close()
			if _, err := store.Load(); err == nil {
				t.Fatal("accepted invalid data")
			}
			after, _ := os.ReadFile(path)
			if string(after) != input {
				t.Fatal("clobbered existing data")
			}
		})
	}
}

func TestFailedSavePreservesLibrary(t *testing.T) {
	path := filepath.Join(t.TempDir(), "library.json")
	store, err := library.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	data, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	data.Lists[0].Entries[0].Lines[0].CW = "UNSUPPORTED!"
	if err := store.Save(data); err == nil {
		t.Fatal("invalid save accepted")
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatal("invalid save replaced data")
	}
	store.Close()
	if err := store.Save(library.Fixture()); err == nil {
		t.Fatal("closed store can save")
	}
	// A directory at the destination causes a real rename failure even as root.
	badPath := filepath.Join(t.TempDir(), "library.json")
	store, err = library.Open(badPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := os.Mkdir(badPath, 0700); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(library.Fixture()); err == nil {
		t.Fatal("rename failure ignored")
	}
	if entries, _ := filepath.Glob(filepath.Join(filepath.Dir(badPath), ".library-*")); len(entries) != 0 {
		t.Fatal("temporary files leaked")
	}
}
