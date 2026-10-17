package library

import (
	"fmt"
	"math"
	"slices"
)

func (s Settings) Validate() error {
	for _, v := range []float64{s.CharacterWPM, s.OverallWPM, s.LearnerHz, s.OtherHz, s.PauseSeconds} {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("settings must be finite")
		}
	}
	if s.CharacterWPM < 5 || s.CharacterWPM > 60 || s.OverallWPM < 5 || s.OverallWPM > s.CharacterWPM {
		return fmt.Errorf("speeds: 5–60 WPM, overall ≤ character")
	}
	if s.LearnerHz < 100 || s.LearnerHz > 2000 || s.OtherHz < 100 || s.OtherHz > 2000 {
		return fmt.Errorf("tones: 100–2000 Hz")
	}
	if s.PauseSeconds < 0 || s.PauseSeconds > 30 {
		return fmt.Errorf("pause: 0–30 seconds")
	}
	if s.SpeechInterval < 1 || s.SpeechInterval > 100 {
		return fmt.Errorf("speech interval: 1–100")
	}
	if !slices.Contains([]string{"cw", "sandwich", "occasional"}, s.Style) || !slices.Contains([]string{"whole", "highlight", "hidden"}, s.Display) || !slices.Contains([]string{"repeat", "cycle", "shuffle"}, s.Sequence) || !slices.Contains([]string{"manual", "automatic"}, s.LineAdvancement) {
		return fmt.Errorf("unknown practice mode")
	}
	return nil
}
