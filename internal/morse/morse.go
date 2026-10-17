// Package morse compiles International Morse into a sample-accurate timeline.
package morse

import (
	"fmt"
	"math"
	"strings"
)

const SampleRate = 48000
const MaxText = 512

var alphabet = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.", 'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..", 'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-", 'Y': "-.--", 'Z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-", '5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
}

type Segment struct {
	Start, End int64
	Tone       bool
}
type Marker struct {
	Start, End int64
	From, To   int
}
type Timeline struct {
	Text     string
	Markers  []Marker
	Segments []Segment
	Samples  int64
}

// Compile ends at the last mark unless calibration requests a trailing word gap.
// Farnsworth spacing follows ARRL's PARIS 50-unit / 19 spacing-unit standard.
func Compile(text string, character, overall float64, calibration bool) (Timeline, error) {
	var result Timeline
	if math.IsNaN(character) || math.IsNaN(overall) || character < 5 || character > 60 || overall < 5 || overall > character {
		return result, fmt.Errorf("speeds must be 5–60 WPM, overall no greater than character")
	}
	if len(text) > MaxText {
		return result, fmt.Errorf("CW text exceeds %d bytes", MaxText)
	}
	text = strings.ToUpper(strings.TrimSpace(text))
	if text == "" {
		return result, fmt.Errorf("CW text is empty")
	}
	tokens, err := Tokens(text)
	if err != nil {
		return result, err
	}
	result.Text = text
	dit := 1.2 / character
	space := (60/overall - 31*dit) / 19
	elapsed := 0.0
	add := func(tone bool, seconds float64) {
		start := result.Samples
		elapsed += seconds
		result.Samples = int64(math.Round(elapsed * SampleRate))
		result.Segments = append(result.Segments, Segment{start, result.Samples, tone})
	}

	for i, token := range tokens {
		if i > 0 {
			gap := 3 * space
			if token.WordGap {
				gap = 7 * space
			}
			add(false, gap)
		}
		start := result.Samples
		for j, e := range token.Code {
			if j > 0 {
				add(false, dit)
			}
			duration := dit
			if e == '-' {
				duration *= 3
			}
			add(true, duration)
		}
		result.Markers = append(result.Markers, Marker{start, result.Samples, token.Start, token.End})
	}

	if calibration {
		add(false, 7*space)
	}
	return result, nil
}
