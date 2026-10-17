package terminal

import (
	tea "charm.land/bubbletea/v2"
	"cw-looper/internal/library"
	"fmt"
	"strconv"
)

func (m *Model) settingsForm() tea.Cmd {
	config := m.data.Lists[m.list].Settings
	labels := []string{"Character WPM", "Overall WPM", "Your tone Hz", "Other tone Hz", "Pause seconds"}
	values := []string{fmt.Sprint(config.CharacterWPM), fmt.Sprint(config.OverallWPM), fmt.Sprint(config.LearnerHz), fmt.Sprint(config.OtherHz), fmt.Sprint(config.PauseSeconds)}
	return m.openForm("List settings · apply on save; speed/tone at next CW", labels, values, func(v []string) (library.Library, error) {
		data := clone(m.data)
		c := &data.Lists[m.list].Settings
		fields := []*float64{&c.CharacterWPM, &c.OverallWPM, &c.LearnerHz, &c.OtherHz, &c.PauseSeconds}
		for i, f := range fields {
			n, err := strconv.ParseFloat(v[i], 64)
			if err != nil {
				return data, fmt.Errorf("%s must be a number", labels[i])
			}
			*f = n
		}
		return data, c.Validate()
	})
}
