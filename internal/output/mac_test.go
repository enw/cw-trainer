//go:build darwin

package output_test

import (
	"context"
	"cw-looper/internal/output"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// Explicit opt-in because these checks play sound on the Mac's selected output.
// Timing here measures API cleanup, not acoustic silence at the speaker.
func TestMacOutputLifecycle(t *testing.T) {
	if os.Getenv("CW_REAL_AUDIO") != "1" {
		t.Skip("set CW_REAL_AUDIO=1 to exercise the Mac output")
	}
	backend := &output.Mac{}
	for _, hz := range []float64{440, 600} {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := backend.CW(ctx, "ET", hz)
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%g Hz ET completed in %s", hz, time.Since(start))
	}
	for _, kind := range []string{"cw", "speech"} {
		t.Run(kind, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() {
				if kind == "cw" {
					done <- backend.CW(ctx, strings.Repeat("KE2ISH ", 20), 440)
				} else {
					done <- backend.Speech(ctx, strings.Repeat("K E two I S H ", 30))
				}
			}()
			time.Sleep(600 * time.Millisecond)
			start := time.Now()
			cancel()
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("wanted cancellation, got %v", err)
				}
				elapsed := time.Since(start)
				t.Logf("%s cancellation/reap/drain returned in %s (not an acoustic measurement)", kind, elapsed)
				if elapsed > time.Second {
					t.Fatal("cleanup exceeded one second")
				}
			case <-time.After(3 * time.Second):
				t.Fatal("output did not cancel")
			}
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := backend.Speech(ctx, "K E two I S H"); err != nil {
		t.Fatal(err)
	}
}
