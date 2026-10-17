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
