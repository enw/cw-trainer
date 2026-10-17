package morse

import (
	"fmt"
	"strings"
)

type Prosign struct{ Token, Code, Meaning string }

var Prosigns = []Prosign{{"<AR>", ".-.-.", "end of message"}, {"<AS>", ".-...", "stand by"}, {"<BT>", "-...-", "break"}, {"<SK>", "...-.-", "end of contact"}, {"<HH>", "........", "correction"}}

type Token struct {
	Text, Code string
	Start, End int
	WordGap    bool
}

func Tokens(text string) ([]Token, error) {
	var result []Token
	wordGap := false
	for i := 0; i < len(text); {
		if text[i] == ' ' {
			wordGap = true
			i++
			continue
		}
		start := i
		code := ""
		if text[i] == '<' {
			end := strings.IndexByte(text[i:], '>')
			if end < 0 {
				return nil, fmt.Errorf("unclosed prosign")
			}
			i += end + 1
			for _, p := range Prosigns {
				if text[start:i] == p.Token {
					code = p.Code
					break
				}
			}
			if code == "" {
				return nil, fmt.Errorf("unknown prosign %q", text[start:i])
			}
		} else {
			code = alphabet[rune(text[i])]
			i++
			if code == "" {
				return nil, fmt.Errorf("unsupported CW character %q", text[start:i])
			}
		}
		result = append(result, Token{text[start:i], code, start, i, wordGap})
		wordGap = false
	}
	return result, nil
}
