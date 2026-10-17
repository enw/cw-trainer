# 01 — Run a saved two-station walking skeleton

**What to build:** Run the real terminal app with one saved two-station exchange: select a line, hear CW, request speech, pause, advance, and quit. This becomes the application foundation.

**Blocked by:** 00 — Design the terminal workflow and playback architecture; design review completed.

**Status:** ready-for-human — implemented; final listening acceptance pending

- [x] Implement the selected minimal stack and document actual install/run/test commands; use real Mac audio, built-in speech, and durable local storage rather than mocked integrations for the demo.
- [x] Load a minimal saved exchange containing KE2ISH and a second station. Preserve line order, station roles, and explicit spoken phrases in the round-tripped model.
- [x] Show the current whole line, station role, playback state, and discoverable controls. Repeat until Next; support selecting/advancing lines with arrows and vi keys as designed.
- [x] Play at fixed 20 WPM character / 15 WPM overall Farnsworth timing, 440 Hz for the learner and 600 Hz for the other station, with two-second pauses.
- [x] S speaks the selected line’s saved phrase. Space, navigation, and exit work during CW, speech, and silence according to the reviewed state design; cancelled output does not overlap or leak into the next line.
- [x] Use standard Morse and PARIS timing verified against authoritative references; retain internal character gaps. The fixture needs no prosign support yet.
- [x] Create initial saved data only when absent, load it after restart, and demonstrate a changed saved fixture survives a save/reload round trip. Report read/save errors without replacing existing data.
- [x] Verify session behavior with controlled outputs and independently verify actual rendered duration/frequency. Test speech arguments as data rather than shell syntax and cancellation at the integration boundary.
- [x] Record the real Mac output lifecycle and PTY walkthrough, including pause/navigation, terminal restoration, and saved data after restart.
- [ ] Confirm audible tone/speech quality and no overlap/tails on the intended output with a listener; the audible stop target is not yet measured.
- [x] Compare software findings with the design and update the same design/drafts.
- [ ] Resolve any failures found during listening acceptance before feature work starts.

## Boundary

This is a narrow end-to-end implementation, not disposable scaffolding. Fixed settings, explicit spoken phrases, one exchange, and whole-line display are sufficient. Full editors, reinforcement patterns, runtime tuning, prosigns, and alternate display modes belong to later tickets.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
User approved the design and authorized this slice. Software verification has passed; final listening acceptance is pending.
The [skeleton fixture and acceptance script](../../../docs/design/skeleton-acceptance.md) define the narrow demo and real Mac checks.

## Result

See [verification evidence](../../../docs/design/skeleton-verification.md) and
[run instructions](../../../README.md). Tests, race checks, vet, real Mac output
lifecycle checks, and a PTY walkthrough passed. No later slice has started.
