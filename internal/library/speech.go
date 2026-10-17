package library

import (
	"cw-looper/internal/morse"
	"strings"
	"unicode"
)

// Spoken follows CW edits until an explicit override is saved.
func Spoken(line Line) string {
	if line.Spoken != "" {
		return line.Spoken
	}
	words := strings.Fields(strings.ToUpper(line.CW))
	for i, word := range words {
		for _, p := range morse.Prosigns {
			word = strings.ReplaceAll(word, p.Token, " "+p.Meaning+" ")
		}
		if word != words[i] {
			words[i] = strings.TrimSpace(word)
			continue
		}
		switch word {
		case "CQ":
			words[i] = "see cue"
		case "73":
			words[i] = "seventy three"
		default:
			abbreviation := strings.Contains(" DE QTH QRS QRQ QRM QRN QSB QRL QRZ QRP QRO QSO QSY RST TU TNX HW CPY BK KN CL R K ES UR FB OM PSE WX ", " "+word+" ")
			hasDigit, hasLetter := false, false
			for _, r := range word {
				hasDigit = hasDigit || unicode.IsDigit(r)
				hasLetter = hasLetter || unicode.IsLetter(r)
			}
			if abbreviation || hasDigit && hasLetter || hasDigit && !hasLetter {
				words[i] = spell(word)
			} else {
				words[i] = strings.ToLower(word)
			}
		}
	}
	return strings.Join(words, " ")
}
func spell(word string) string {
	digits := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	parts := make([]string, 0, len(word))
	for _, r := range word {
		if r >= '0' && r <= '9' {
			parts = append(parts, digits[r-'0'])
		} else {
			parts = append(parts, string(r))
		}
	}
	return strings.Join(parts, " ")
}
