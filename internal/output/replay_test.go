package output_test

import (
	"context"
	"cw-looper/internal/library"
	"cw-looper/internal/morse"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"testing"
)

type replayBackend struct {
	timeline morse.Timeline
	hz       float64
}

func (b *replayBackend) CW(context.Context, string, float64, float64, float64) error { return nil }
func (b *replayBackend) Speech(context.Context, string) error                        { return nil }
func (b *replayBackend) Timeline(_ context.Context, timeline morse.Timeline, hz float64) error {
	b.timeline = timeline
	b.hz = hz
	return nil
}
func TestMeasuredReplayUsesSelectedStationTone(t *testing.T) {
	backend := &replayBackend{}
	runner := output.NewRunner(backend)
	defer runner.Close()
	timeline := morse.Timeline{Samples: 3888, Segments: []morse.Segment{{Start: 0, End: 3888, Tone: true}}}
	settings := library.Defaults()
	settings.OtherHz = 750
	event := <-runner.Execute(session.Effect{Kind: session.PlayAttempt, ID: 9, Timeline: timeline, Settings: settings}, library.Line{Station: "other"})
	if event.Kind != session.Done || backend.hz != 750 || backend.timeline.Samples != 3888 {
		t.Fatal("replay replaced measured timing or station tone")
	}
}
