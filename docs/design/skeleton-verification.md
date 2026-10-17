# Walking-skeleton verification

2026-10-16. Ticket 01 implementation on macOS 26.6.2 / Apple silicon, Go 1.26.5.
Pinned direct dependencies: Bubble Tea v2.1.0 and Oto v3.5.1. Design approved by
the user before implementation. No later feature slice has started.

## Implemented

One seeded, saved two-station exchange; silent list/exchange/line browsing;
whole-line practice; repeat-until-Next; fixed 20/15 WPM and 440/600 Hz; two-second
gaps; pause/resume; speak-now; arrow/vi navigation; stop/quit cleanup.

The session is the single playback policy owner. Workers acknowledge stop only
after cancellation and cleanup, and operation IDs remain unique across exchange
selection. A regression check covers top-level navigation while paused.

## Automated results

- `go test ./...`: passed.
- `go test -race ./...`: passed.
- `go vet ./...`: passed.
- `go build -o bin/cw-looper ./cmd/cw-looper`: passed.
- Independent Morse cases verify PARIS calibration, mark/gap lengths, PCM
  duration, 600 Hz cycle count, envelope amplitude, and seeking consistency.
- Store checks cover seed/load/edit/reload, empty libraries, malformed/future
  data preservation, exclusive locking, invalid/closed saves, actual rename
  failure, and temporary-file cleanup.
- Session/terminal checks cover cancellation barriers, stale completions,
  repetition counting, pause/restart, frozen gap remainder, paused speech and
  navigation, speech interruption, failure handling, quit priority, and exchange
  generation uniqueness.

## Real Mac output lifecycle

`CW_REAL_AUDIO=1 go test -v ./internal/output -run TestMacOutputLifecycle -count=1`
passed outside the execution sandbox. The sandboxed attempt failed at
AudioQueueStart; unsandboxed execution reached the real audio backend successfully.

| Check | Observed |
| --- | --- |
| 440 Hz ET, including initial audio setup and handoff allowance | 745.6 ms |
| 600 Hz ET, including handoff allowance | 626.2 ms |
| CW cancellation, software cleanup and handoff allowance | 51.1 ms |
| Speech cancellation, child reaping and handoff allowance | 52.0 ms |
| Spoken KE2ISH phrase | Process completed successfully |

These are wall-clock API observations, not acoustic measurements. A 50 ms output
handoff allowance is implemented after pause/completion/reaping. It is not a
claim that all output devices drain within that time. Built-in/wired audible
tail, pronunciation, and click quality still require listening validation;
Bluetooth/device changes are unverified.

## Actual terminal walkthrough

Used the built binary in a 100 × 28 PTY with real Mac outputs and an isolated
library under `/tmp`, leaving the user's normal library untouched. Observed:

- Enter through list/exchange/line into 600 Hz Other practice.
- Space stops into Paused; l selects the learner line silently.
- s enters Speaking, then returns to Paused.
- Space restarts CW at 440 Hz; s interrupts it for speech.
- h/l navigation during speech cancels before starting the new line.
- q during speech exits; `stty -a` afterward reports canonical input and echo
  restored, with the alternate screen exited and cursor restored.
- After editing the closed app's temporary JSON, relaunch shows the changed
  `KE2ISH` line. First-run seeding did not replace it.

The PTY output was inspected as terminal text/control sequences, not as a native
terminal screenshot. No claim of rendered visual QA or human listening is made.
Pause-during-gap behavior is covered with controlled time; its exact live timing
has not been instrumented acoustically.

## Design refinements

- Skeleton JSON accepts only exchanges with explicit spoken phrases and the
  fixed preset. It rejects unsupported modes/settings rather than ignoring them.
  Ordinary items and generated speech remain later slices.
- Input length is bounded: CW 512 bytes, spoken phrase 1024 bytes, library 2 MiB,
  up to 100 lists/entries/lines at each level. PCM is streamed, not preallocated
  for full phrases. Storage errors are reported before terminal startup.
- Audio/device errors are recoverable by retry or restart as indicated. There is
  a playback watchdog; no retry loop runs automatically.
- Character progress estimation is deferred to ticket 05; the skeleton displays
  the full line and therefore has no competing visual timing scheduler.

## Acceptance still requiring a listener

Run the app on the intended Mac output and confirm both tones, clear CW without
clicks, recognizable speech, and no audible overlap/tail during pause, speak-now,
navigation, and quit. The original ≤100 ms audible-stop target remains unverified.
If that check fails, revise the output adapter before advancing feature work.
Ticket 01 is implemented and ready for this final acceptance, not marked resolved.
