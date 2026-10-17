# 01 — Run a saved two-station walking skeleton

**What to build:** Run the real terminal app with one saved two-station exchange: select a line, hear CW, request speech, pause, advance, and quit. This becomes the application foundation.

**Blocked by:** 00 — Design the terminal workflow and playback architecture; design review completed.

**Status:** claimed — implementation and verification in progress

- [ ] Implement the selected minimal stack and document actual install/run/test commands; use real Mac audio, built-in speech, and durable local storage rather than mocked integrations for the demo.
- [ ] Load a minimal saved exchange containing KE2ISH and a second station. Preserve line order, station roles, and explicit spoken phrases in the round-tripped model.
- [ ] Show the current whole line, station role, playback state, and discoverable controls. Repeat until Next; support selecting/advancing lines with arrows and vi keys as designed.
- [ ] Play at fixed 20 WPM character / 15 WPM overall Farnsworth timing, 440 Hz for the learner and 600 Hz for the other station, with two-second pauses.
- [ ] S speaks the selected line’s saved phrase. Space, navigation, and exit work during CW, speech, and silence according to the reviewed state design; cancelled output does not overlap or leak into the next line.
- [ ] Use standard Morse and PARIS timing verified against authoritative references; retain internal character gaps. The fixture needs no prosign support yet.
- [ ] Create initial saved data only when absent, load it after restart, and demonstrate a changed saved fixture survives a save/reload round trip. Report read/save errors without replacing existing data.
- [ ] Verify session behavior with controlled outputs and independently verify actual rendered duration/frequency. Test speech arguments as data rather than shell syntax and cancellation at the integration boundary.
- [ ] Perform and record a real Mac walkthrough demonstrating responsive input, both tones, actual speech, pause/navigation, terminal restoration, and saved data after restart.
- [ ] Compare findings with the design, update the same design/drafts, and resolve material stack or playback failures before feature work starts.

## Boundary

This is a narrow end-to-end implementation, not disposable scaffolding. Fixed settings, explicit spoken phrases, one exchange, and whole-line display are sufficient. Full editors, reinforcement patterns, runtime tuning, prosigns, and alternate display modes belong to later tickets.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
User approved the design and authorized this slice. Skeleton validation is in progress.
The [skeleton fixture and acceptance script](../../../docs/design/skeleton-acceptance.md) define the narrow demo and real Mac checks.
