// Package output runs cancellable effects; it does not decide practice policy.
package output

import (
	"context"
	"cw-looper/internal/library"
	"cw-looper/internal/session"
	"fmt"
	"sync"
	"time"
)

type Backend interface {
	CW(context.Context, string, float64, float64, float64) error
	Speech(context.Context, string) error
}
type operation struct {
	id     uint64
	cancel context.CancelFunc
	done   chan struct{}
}
type Runner struct {
	backend Backend
	mu      sync.Mutex
	active  *operation
}

func NewRunner(backend Backend) *Runner { return &Runner{backend: backend} }

// Execute returns immediately. A stop acknowledgment waits for backend cleanup,
// including process reaping/drain, even if a completion is already in flight.
func (r *Runner) Execute(effect session.Effect, line library.Line) <-chan session.Input {
	result := make(chan session.Input, 1)
	r.mu.Lock()
	defer r.mu.Unlock()
	if effect.Kind == session.StopOutput {
		op := r.active
		if op != nil && op.id == effect.ID {
			op.cancel()
		} else {
			op = nil
		}
		go func() {
			if op != nil {
				<-op.done
			}
			result <- session.Input{Kind: session.Stopped, ID: effect.ID}
			close(result)
		}()
		return result
	}
	if r.active != nil {
		select {
		case <-r.active.done:
		default:
			result <- session.Input{Kind: session.Error, ID: effect.ID, Err: fmt.Errorf("previous output has not stopped")}
			close(result)
			return result
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	op := &operation{effect.ID, cancel, make(chan struct{})}
	r.active = op
	go func() {
		defer cancel()
		var err error
		switch effect.Kind {
		case session.PlayCW:
			config := effect.Settings
			if config.CharacterWPM == 0 {
				config = library.Defaults()
			}
			hz := config.LearnerHz
			if line.Station == "other" {
				hz = config.OtherHz
			}
			err = r.backend.CW(ctx, line.CW, hz, config.CharacterWPM, config.OverallWPM)
		case session.PlaySpeech:
			err = r.backend.Speech(ctx, line.Spoken)
		case session.Wait:
			timer := time.NewTimer(effect.Duration)
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case <-timer.C:
			}
			timer.Stop()
		default:
			err = fmt.Errorf("unsupported output effect")
		}
		close(op.done)
		event := session.Input{Kind: session.Done, ID: effect.ID}
		if err != nil {
			event.Kind = session.Error
			event.Err = err
		}
		result <- event
		close(result)
	}()
	return result
}
func (r *Runner) Close() {
	r.mu.Lock()
	op := r.active
	if op != nil {
		op.cancel()
	}
	r.mu.Unlock()
	if op != nil {
		<-op.done
	}
}
