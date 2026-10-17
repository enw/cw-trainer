// Package library owns the saved practice material and its validation.
package library

import (
	"cw-looper/internal/morse"
	"fmt"
	"strings"
	"unicode"
)

type Library struct {
	SchemaVersion int    `json:"schemaVersion"`
	SeedVersion   int    `json:"seedVersion"`
	Lists         []List `json:"lists"`
}
type List struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Entries  []Entry  `json:"entries"`
	Settings Settings `json:"settings"`
}
type Entry struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Title string `json:"title"`
	Lines []Line `json:"lines"`
}
type Line struct {
	ID      string `json:"id"`
	Station string `json:"station"`
	CW      string `json:"cw"`
	Spoken  string `json:"spokenOverride"`
}
type Settings struct {
	CharacterWPM    float64 `json:"characterWPM"`
	OverallWPM      float64 `json:"overallWPM"`
	LearnerHz       float64 `json:"learnerHz"`
	OtherHz         float64 `json:"otherHz"`
	PauseSeconds    float64 `json:"pauseSeconds"`
	Style           string  `json:"style"`
	Display         string  `json:"display"`
	Sequence        string  `json:"sequence"`
	LineAdvancement string  `json:"lineAdvancement"`
	SpeechInterval  int     `json:"speechInterval"`
}

func Defaults() Settings { return Settings{20, 15, 440, 600, 2, "cw", "whole", "repeat", "manual", 5} }
func Fixture() Library {
	return Library{SchemaVersion: 1, SeedVersion: 1, Lists: []List{{
		ID: "starter", Name: "First exchange", Settings: Defaults(), Entries: []Entry{{
			ID: "call-answer", Kind: "exchange", Title: "Call and answer", Lines: []Line{
				{"calling", "other", "CQ CQ DE W1AW K", "see cue see cue D E W one A W K"},
				{"answering", "learner", "W1AW DE KE2ISH K", "W one A W D E K E two I S H K"},
			},
		}},
	}}}
}
func (l Library) Validate() error {
	if l.SchemaVersion != 1 {
		return fmt.Errorf("unsupported schema version %d (expected 1)", l.SchemaVersion)
	}
	if l.SeedVersion < 0 {
		return fmt.Errorf("invalid seed version")
	}
	ids := map[string]bool{}
	id := func(s string) error {
		if !plain(s, 100) || ids[s] {
			return fmt.Errorf("IDs must be unique nonempty plain text")
		}
		ids[s] = true
		return nil
	}
	if len(l.Lists) > 100 {
		return fmt.Errorf("too many lists (maximum 100)")
	}
	for _, list := range l.Lists {
		if err := id(list.ID); err != nil {
			return err
		}
		if !plain(list.Name, 100) {
			return fmt.Errorf("invalid list name")
		}
		if err := list.Settings.Validate(); err != nil {
			return err
		}
		if len(list.Entries) > 100 {
			return fmt.Errorf("too many entries")
		}
		for _, entry := range list.Entries {
			if err := id(entry.ID); err != nil {
				return err
			}
			if (entry.Kind != "exchange" && entry.Kind != "item") || !plain(entry.Title, 100) {
				return fmt.Errorf("expected an item or exchange with a title")
			}
			if len(entry.Lines) == 0 || len(entry.Lines) > 100 || (entry.Kind == "item" && len(entry.Lines) != 1) {
				return fmt.Errorf("exchange must contain 1–100 lines")
			}
			for _, line := range entry.Lines {
				if err := id(line.ID); err != nil {
					return err
				}
				if line.Station != "learner" && line.Station != "other" {
					return fmt.Errorf("invalid station role")
				}
				if !plain(line.CW, morse.MaxText) {
					return fmt.Errorf("CW text must be nonempty plain text")
				}
				if _, err := morse.Compile(line.CW, 20, 15, false); err != nil {
					return err
				}
				if line.Spoken != "" && !plain(line.Spoken, 1024) {
					return fmt.Errorf("spoken phrase must contain 1–1024 plain characters")
				}
			}
		}
	}
	return nil
}
func plain(s string, limit int) bool {
	if strings.TrimSpace(s) == "" || len(s) > limit {
		return false
	}
	for _, c := range s {
		if unicode.IsControl(c) {
			return false
		}
	}
	return true
}
