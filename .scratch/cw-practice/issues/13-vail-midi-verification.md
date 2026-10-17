# 13 — Verify Vail touch-paddle input over USB MIDI

**What to build:** Select the learner's Vail adapter on the Mac and see timestamped key-down/key-up events while using its capacitive touch pads, proving the input path before building exercises.

**Blocked by:** None — independent hardware integration check; does not depend on 11/12.

**Status:** hardware input verified — app diagnostic and robustness work remain queued

- [ ] Identify the actual adapter board/firmware where possible and enumerate Mac MIDI endpoints. Distinguish an enumerated endpoint from verified keying events.
- [ ] Provide a terminal diagnostic that selects the adapter, shows connection state and event timestamps, and offers a deliberate MIDI-mode connection action.
- [ ] Verify the documented mode switch on this device. Preserve its existing Iambic B keyer, speed and sidetone; do not switch to passthrough or rewrite persistent settings just to test input.
- [x] With the user touching the pads, capture isolated dits/dahs, held pads, and a squeeze/release. Confirm balanced down/up events and measured mark/gap durations. Record actual observations separately from protocol expectations.
- [ ] Treat Iambic B output as already-timed marks; do not assume MIDI note numbers identify dits versus dahs in this mode. Account for Note Off and conventional zero-velocity Note On release handling.
- [x] Explain that MIDI carries event timing, not microphone audio. The adapter provides live sidetone; hearing its buzzer alone is not proof that the Mac receives events.
- [ ] Handle missing devices, silent endpoints, disconnect during a mark, reconnect and leaving the diagnostic without stuck input or interference with normal practice.
- [ ] Document mode ownership and cleanup: explicitly choose the return mode rather than pretending the original mode can be queried/restored. Recheck mode after reconnect and do not rely on it surviving a power cycle.
- [ ] Establish the smallest tested macOS MIDI input boundary, timestamp source and latency limitations for ticket 14. Use a recorded event fixture for automated checks; require a real-device touch test to complete hardware acceptance.

## Known setup and evidence boundary

The learner owns a Vail adapter, uses its capacitive pads without external
paddles, and has selected Iambic B. A cooperative standalone CoreMIDI probe verified this device's MIDI marks on
2026-10-16. MIDI mode was temporarily enabled and keyboard output restored.
No firmware, keyer, speed or tone changes were made; the app diagnostic remains
unimplemented. See [observed evidence](../../../docs/design/vail-midi-verification.md).

[Upstream MIDI specification](https://github.com/Vail-CW/vail-adapter/blob/master/docs/MIDI_INTEGRATION_SPEC.md)
currently documents keyboard mode at boot, explicit MIDI-mode selection using
CC0, and note 0 on/off for onboard-keyer output. The actual device produced notes 1/2 with regular timed repetitions. Support
that observed behavior; note numbers alone do not establish passthrough. Exact
firmware and independent keyer-mode identification remain unverified.
