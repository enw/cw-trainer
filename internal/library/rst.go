package library

import "fmt"

// RatingScale paraphrases ARRL's RST descriptions for a compact reference.
type RatingScale struct {
	Name     string
	Meanings []string
}

func RSTScales() []RatingScale {
	return []RatingScale{
		{"Readability", []string{"Nothing intelligible", "Only isolated words understood", "Copy requires substantial effort", "Copy requires little effort", "Everything clearly understood"}},
		{"Strength", []string{"Just detectable", "Very low level", "Low level", "Fair level", "Fairly good level", "Good level", "Moderately high level", "High level", "Exceptionally high level"}},
		{"Tone", []string{"Broad, coarse AC note at 60 Hz or below", "Broad, harsh AC note", "Unfiltered rectified AC note", "Rough note with slight filtering", "Filtered note with heavy ripple", "Filtered note with noticeable ripple", "Almost pure note with residual ripple", "Almost perfect note with slight modulation", "Pure note without ripple or modulation"}},
	}
}

// ReportSpeech validates three numeric CW ratings; voice reports omit tone.
func ReportSpeech(report string) (string, error) {
	if len(report) != 3 {
		return "", fmt.Errorf("CW RST needs three ratings")
	}
	scales := RSTScales()
	speech := "R S T"
	for i, scale := range scales {
		rating := int(report[i] - '0')
		if report[i] < '1' || rating > len(scale.Meanings) {
			return "", fmt.Errorf("%s rating out of range", scale.Name)
		}
		speech += fmt.Sprintf(". %s %d: %s", scale.Name, rating, scale.Meanings[rating-1])
	}
	return speech, nil
}

func RSTStarter() List {
	var words [][2]string
	for _, report := range []string{"599", "579", "339", "449", "567"} {
		speech, _ := ReportSpeech(report)
		words = append(words, [2]string{report, speech})
	}
	list := starterItems("rst", "RST reports · readability strength tone", words)
	list.Settings.Style = "sandwich"
	first, _ := ReportSpeech("579")
	second, _ := ReportSpeech("339")
	list.Entries = append(list.Entries, Entry{ID: "rst-exchange", Kind: "exchange", Title: "Exchanging reports · practice example", Lines: []Line{
		{ID: "rst-exchange-learner", Station: "learner", CW: "W1AW DE KE2ISH UR RST 579 K", Spoken: "KE2ISH gives W1AW a report. " + first},
		{ID: "rst-exchange-other", Station: "other", CW: "KE2ISH DE W1AW R UR RST 339 K", Spoken: "W1AW gives KE2ISH a report. " + second},
	}})
	return list
}
