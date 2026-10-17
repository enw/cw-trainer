package library

import (
	"cw-looper/internal/morse"
	"strings"
	"unicode"
)

// Spoken uses the same tokens as CW, including adjacent prosigns. Ordinary
// spans are pronounced together so words stay natural and callsigns are spelled.
func Spoken(line Line) string {
	if line.Spoken != "" {
		return line.Spoken
	}
	tokens, err := morse.Tokens(strings.ToUpper(strings.TrimSpace(line.CW)))
	if err != nil {
		return line.CW
	} // Validation normally rejects unsupported text.
	var result []string
	var word strings.Builder
	flush := func() {
		if word.Len() > 0 {
			result = append(result, pronounce(word.String()))
			word.Reset()
		}
	}
	for _, token := range tokens {
		if token.WordGap {
			flush()
		}
		if strings.HasPrefix(token.Text, "<") {
			flush()
			for _, p := range morse.Prosigns {
				if p.Token == token.Text {
					result = append(result, p.Meaning)
					break
				}
			}
		} else {
			word.WriteString(token.Text)
		}
	}
	flush()
	return strings.Join(result, " ")
}
func pronounce(word string) string {
	switch word {
	case "CQ":
		return "see cue"
	case "73":
		return "seventy three"
	case "AGN":
		return "again"
	case "PSE":
		return "please"
	case "FER":
		return "for"
	}
	abbreviation := strings.Contains(" DE QTH QRS QRQ QRM QRN QSB QRL QRZ QRP QRO QSO QSY RST TU TNX HW CPY BK KN CL R K ES UR FB OM WX ", " "+word+" ")
	hasDigit := false
	for _, r := range word {
		hasDigit = hasDigit || unicode.IsDigit(r)
	}
	if abbreviation || hasDigit {
		return spell(word)
	}
	return strings.ToLower(word)
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
