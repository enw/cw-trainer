package output_test

import (
	"context"
	"cw-looper/internal/library"
	"cw-looper/internal/output"
	"cw-looper/internal/session"
	"testing"
	"time"
)

type controlledOutput struct {
	started chan string
	release chan struct{}
}

func (f *controlledOutput) CW(ctx context.Context, text string, hz, character, overall float64) error {
	return f.play(ctx, "cw")
}
func (f *controlledOutput) Speech(ctx context.Context, text string) error { return f.play(ctx, text) }
func (f *controlledOutput) play(ctx context.Context, text string) error {
	f.started <- text
	<-ctx.Done()
	<-f.release
	return ctx.Err()
}
func TestStopAcknowledgesOnlyAfterOutputDrains(t *testing.T) {
	fake := &controlledOutput{make(chan string, 2), make(chan struct{})}
	runner := output.NewRunner(fake)
	line := library.Line{CW: "KE2ISH", Station: "learner", Spoken: "K E two I S H"}
	done := runner.Execute(session.Effect{Kind: session.PlayCW, ID: 1}, line)
	select {
	case <-fake.started:
	case <-time.After(time.Second):
		t.Fatal("not started")
	}
	stopped := runner.Execute(session.Effect{Kind: session.StopOutput, ID: 1}, line)
	select {
	case <-stopped:
		t.Fatal("acknowledged before drain")
	case <-time.After(10 * time.Millisecond):
	}
	close(fake.release)
	if event := <-stopped; event.Kind != session.Stopped || event.ID != 1 {
		t.Fatal(event)
	}
	<-done
	speech := runner.Execute(session.Effect{Kind: session.PlaySpeech, ID: 2}, line)
	if text := <-fake.started; text != line.Spoken {
		t.Fatal("speech not passed as data")
	}
	runner.Close()
	<-speech
}
