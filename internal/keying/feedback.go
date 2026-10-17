package keying

import (
	"cw-looper/internal/morse"
	"errors"
	"fmt"
	"math"
	"strings"
)

// Feedback is a tentative interpretation of measured timing, never of a target.
type Feedback struct {
	Text, Status, Timing string
	Uncertain            bool
	groups               []morse.Token
}

// Feedback uses sending WPM, independently of reference or Farnsworth settings.
func (a Attempt) Feedback(wpm float64) Feedback {
	f := Feedback{}
	if math.IsNaN(wpm) || math.IsInf(wpm, 0) || wpm < 5 || wpm > 60 {
		f.Status = "Set sending speed between 5 and 60 WPM"
		return f
	}
	if a.Err != nil {
		f.Status = "Incomplete capture: " + a.Err.Error()
		if errors.Is(a.Err, ErrNoInput) {
			f.Status = "No input"
		}
		return f
	}
	if a.Active || !a.complete {
		f.Status = "Capture not complete"
		if !a.Active && len(a.Marks) == 0 {
			f.Status = "No input"
		}
		return f
	}
	unit := 1200 / wpm
	code, word := "", false
	var text strings.Builder
	gaps := map[string][]float64{}
	flush := func() {
		if code == "" {
			return
		}
		if len(f.groups) > 0 {
			if word {
				text.WriteString(" ")
			}
		}
		label := patternLabel(code)
		if label == "" {
			label = "[" + code + "?]"
			f.Uncertain = true
		}
		text.WriteString(label)
		f.groups = append(f.groups, morse.Token{Code: code, WordGap: word})
		code = ""
	}
	for i, mark := range a.Marks {
		if i > 0 {
			gap := float64(mark.Start-a.Marks[i-1].End) / 1e6 / unit
			kind := classify(gap)
			switch kind {
			case "dit":
				gaps["element"] = append(gaps["element"], gap)
			case "dah":
				flush()
				word = false
				gaps["character"] = append(gaps["character"], gap)
			case "word":
				flush()
				word = true
				gaps["word"] = append(gaps["word"], gap)
			default:
				flush()
				text.WriteString("[gap?]")
				word = false
				f.Uncertain = true
				gaps["uncertain"] = append(gaps["uncertain"], gap)
			}
		}
		units := float64(mark.End-mark.Start) / 1e6 / unit
		switch classify(units) {
		case "dit":
			code += "."
		case "dah":
			code += "-"
		default:
			code += "?"
			f.Uncertain = true
		}
	}
	flush()
	f.Text = text.String()
	f.Status = "Tentative decode"
	if f.Uncertain {
		f.Status = "Uncertain timing or symbol"
	}
	var ranges []string
	for _, kind := range []string{"element", "character", "word", "uncertain"} {
		values := gaps[kind]
		if len(values) == 0 {
			continue
		}
		lo, hi := values[0], values[0]
		for _, v := range values {
			lo = min(lo, v)
			hi = max(hi, v)
		}
		ranges = append(ranges, fmt.Sprintf("%s %.1f–%.1f", kind, lo, hi))
	}
	f.Timing = "Gaps (dit units): " + strings.Join(ranges, "; ")
	if len(ranges) == 0 {
		f.Timing = "No internal gaps to assess"
	}
	return f
}

func classify(units float64) string {
	switch {
	case units >= 0.65 && units <= 1.35:
		return "dit"
	case units >= 2.4 && units <= 3.6:
		return "dah"
	case units >= 5.6 && units <= 8.4:
		return "word"
	default:
		return "uncertain"
	}
}

func patternLabel(code string) string {
	var labels []string
	for _, char := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" {
		tokens, _ := morse.Tokens(string(char))
		if tokens[0].Code == code {
			labels = append(labels, string(char))
		}
	}
	for _, p := range morse.Prosigns {
		if p.Code == code {
			labels = append(labels, p.Token)
		}
	}
	// These supported prosign patterns also denote punctuation on air.
	switch code {
	case ".-.-.":
		labels = append(labels, "+")
	case "-...-":
		labels = append(labels, "=")
	}
	return strings.Join(labels, "/")
}

// Compare checks code groups and word boundaries only after independent decode.
func (f Feedback) Compare(target string) string {
	if len(f.groups) == 0 {
		return f.Status
	}
	if f.Uncertain {
		return "Comparison uncertain; adjust sending WPM or replay spacing"
	}
	tokens, err := morse.Tokens(strings.ToUpper(strings.TrimSpace(target)))
	if err != nil {
		return "Target cannot be compared"
	}
	match := len(tokens) == len(f.groups)
	if match {
		for i, t := range tokens {
			if t.Code != f.groups[i].Code || i > 0 && t.WordGap != f.groups[i].WordGap {
				match = false
				break
			}
		}
	}
	if match {
		return "Matching Morse patterns and word boundaries"
	}
	return "Different Morse patterns or word boundaries"
}
