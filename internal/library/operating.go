package library

import "strings"

// PhoneticWords follows the ITU spelling alphabet, including Alfa and Juliett.
// Source: https://www.arrl.org/quick-reference-operating-aids
func PhoneticWords() []string {
	return strings.Fields("Alfa Bravo Charlie Delta Echo Foxtrot Golf Hotel India Juliett Kilo Lima Mike November Oscar Papa Quebec Romeo Sierra Tango Uniform Victor Whiskey X-ray Yankee Zulu")
}

func PhoneticStarter() List {
	var words [][2]string
	for i, word := range PhoneticWords() {
		words = append(words, [2]string{string(rune('A' + i)), word})
	}
	list := starterItems("phonetic", "NATO / ITU phonetic alphabet", words)
	list.Settings.Style = "sandwich"
	return list
}
