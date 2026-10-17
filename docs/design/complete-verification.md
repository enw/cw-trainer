# Complete implementation verification — 2026-10-16

Tickets 02–10 were implemented in dependency order on `complete-cw-practice`,
with one commit per feature slice. Tickets 11–12 are later requested drafts,
not implemented scope. The original skeleton record remains historical.

## Automated evidence

Go 1.26.5, macOS 26.6.2, Apple silicon:

- `go test ./...`: passed.
- `go test -race ./...`: passed.
- `go vet ./...`: passed.
- `go build -o bin/cw-looper ./cmd/cw-looper`: passed.
- Module metadata tidied; build output is ignored by Git.

Coverage includes independent PARIS/Farnsworth timing and sampled PCM,
continuous prosigns and marker spans, stale progress/completion rejection,
output cleanup barriers, runtime timing snapshots, all reinforcement styles,
unit-boundary sequencing, deterministic shuffle rounds, empty/single-item
behavior, paused navigation, and ordered exchange progression. The complete
style × sequence × line-mode matrix is tested at the session boundary.

Terminal/store checks cover list/item creation, typed shortcuts and paste focus,
exchange role/speech editing and reorder/save/reload, hidden-title concealment,
whole-prosign highlighting, long ordinary phrases, save/rename failures,
locking, malformed/future-version data preservation, starter upgrades, retained
per-list settings, and deliberately empty libraries/lists.

## Real Mac output

`CW_REAL_AUDIO=1 go test -v ./internal/output -run TestMacOutputLifecycle -count=1`
passed using actual Oto and macOS `say`:

- 440 Hz ET: 741.6 ms including cold initialization and drain allowance.
- 600 Hz ET: 621.2 ms including drain allowance.
- CW cancellation cleanup: 50.8 ms.
- Speech cancellation/reap cleanup: 51.6 ms.
- Callsign speech completed successfully.

These are API/lifecycle observations, not acoustic measurements. The 50 ms
handoff allowance remains a conservative software policy, not proven silence.

## Interactive PTY observations

The actual built application ran at 80 × 24 using an isolated `/tmp` library:

1. First launch showed all four starter lists. Enter opened callsign practice
   with the initial 20/15 WPM, 440 Hz and repeat/CW-only defaults.
2. Space paused; p opened the real settings form. Tab navigation and editing
   saved 60/60 WPM, zero pause, hidden display and CW-talk-CW cycle practice.
   R exposed the answer; Space resumed. Output status progressed CW → Speaking
   → CW, then selected another item and hid its answer again.
3. Speak-now, pause and Esc returned to browsing. q restored the normal terminal
   screen/cursor and exited with code zero.
4. The binary was rebuilt/replaced and restarted with the same data path. Saved
   timing, hidden display and reinforcement/sequence settings were retained.
5. Changing to highlighting, occasional speech every one CW and shuffle showed
   character brackets moving through the actual callsign PCM and speech phases.
6. The exchange list was set to automatic lines at 60/60. The first three-line
   exchange progressed in order, showing Other/600, You/440, Other/600. Pause,
   arrow navigation and quit worked; terminal restoration was observed again.

The first walkthrough exposed clipped footer text at 80 columns. Controls were
split into shorter lines and checked in the rebuilt app. Paste/focus, long-text,
reveal and exchange-edit correctness also have automated regressions.

## Review and fixes

Separate standards and spec reviews compared `7091bab...HEAD`.
Standards found no hard violations and two consequential consistency concerns;
spec found three behavior defects. Their overlapping findings were resolved:

- Ordinary item previews are bounded independently of full 512-byte CW content.
- Speech reuses the CW tokenizer, preserving callsign/CQ pronunciation beside
  adjacent prosigns.
- A pending save blocks concurrent edits but leaves practice controls responsive.
  Normal exit stops playback then waits for that save; emergency exit reports
  unsaved state. A still-running emergency save retains its OS lock until exit.

## Remaining human acceptance

The agent did not listen to speaker output. Before treating the audio experience
as accepted, a person should check:

- Whole-call/word recognition quality, smooth element envelopes, no objectionable
  clicks, and a comfortable distinction between both station tones.
- Actual pronunciation of generated phrases and chosen spoken overrides.
- Highlight alignment with what is heard. Position is an estimate: PCM read
  samples minus queued player bytes and the configured 20 ms device buffer,
  polled by the UI at 25 ms; Bluetooth/external-device latency is not measured.
- Audible stop tails during CW, speech, navigation and rapid interruption; the
  intended ≤100 ms acoustic stop target remains unverified.

This is a listening/visual acceptance boundary, not deferred feature code.
No platform other than this Mac setup or mid-session audio-device switching was
validated.
