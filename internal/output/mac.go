//go:build darwin

package output

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"cw-looper/internal/morse"
	"github.com/ebitengine/oto/v3"
)

// DrainTail is a conservative software handoff allowance, not a measured
// hardware guarantee. Real-device evidence is documented separately.
const DrainTail = 50 * time.Millisecond

type Mac struct {
	once    sync.Once
	audio   *oto.Context
	ready   chan struct{}
	initErr error
}

func (m *Mac) prepare(ctx context.Context) error {
	m.once.Do(func() {
		m.audio, m.ready, m.initErr = oto.NewContext(&oto.NewContextOptions{
			SampleRate: morse.SampleRate, ChannelCount: 1, Format: oto.FormatFloat32LE, BufferSize: 20 * time.Millisecond,
		})
	})
	if m.initErr != nil {
		return fmt.Errorf("audio initialization failed; restart app: %w", m.initErr)
	}
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.ready:
	case <-timer.C:
		return fmt.Errorf("audio initialization timed out; check selected output and restart")
	}
	if err := m.audio.Err(); err != nil {
		return fmt.Errorf("audio device failed; restart app: %w", err)
	}
	return nil
}
func (m *Mac) CW(ctx context.Context, text string, hz, character, overall float64) error {
	return m.CWProgress(ctx, text, hz, character, overall, func(int64) {})
}

func (m *Mac) CWProgress(ctx context.Context, text string, hz, character, overall float64, progress func(int64)) error {
	timeline, err := morse.Compile(text, character, overall, false)
	if err != nil {
		return err
	}
	pcm, err := morse.NewPCM(timeline, hz)
	if err != nil {
		return err
	}
	if err := m.prepare(ctx); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	player := m.audio.NewPlayer(pcm)
	player.SetBufferSize(morse.SampleRate * 4 * 20 / 1000)
	player.Play()
	defer func() {
		player.PauseAndStopReading()
		// Audio already handed to the OS cannot be withdrawn by pausing a player.
		time.Sleep(DrainTail)
	}()
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	// A non-progressing driver must fail rather than hold a session indefinitely.
	deadline := time.NewTimer(time.Duration(float64(timeline.Samples)/morse.SampleRate*float64(time.Second)) + 5*time.Second)
	defer deadline.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("audio playback stalled; check selected output and restart")
		case <-ticker.C:
			// Subtract queued player bytes and the configured device buffer.
			progress(max(0, pcm.ReadSamples()-int64(player.BufferedSize()/4)-morse.SampleRate*20/1000))
			if err := m.audio.Err(); err != nil {
				return fmt.Errorf("audio device failed; restart app: %w", err)
			}
			if err := player.Err(); err != nil {
				return fmt.Errorf("audio playback failed: %w", err)
			}
			if !player.IsPlaying() {
				return nil
			}
		}
	}
}
func (m *Mac) Speech(ctx context.Context, text string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "/usr/bin/say", "-f", "-")
	cmd.Stdin = strings.NewReader(text)
	cmd.WaitDelay = time.Second
	err := cmd.Run()
	// Run waits/reaps the child. Give any submitted device audio the same handoff
	// allowance as CW; this still requires listening validation on the real output.
	time.Sleep(DrainTail)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("speech failed; check the system voice or use CW only: %w", err)
	}
	return nil
}
