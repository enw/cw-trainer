# 14 — Hear a phrase, send it back, and replay the attempt

**What to build:** Play a saved practice item or exchange line, let the learner send it on the Vail touch pads, then replay the captured attempt alongside the reference.

**Blocked by:** 13 — Verified Vail MIDI input; 10 — Existing integrated practice app (implemented).

**Status:** draft — queued follow-up; implementation not requested

- [ ] Add an explicit listen/send-back practice mode using existing saved items and exchange lines, without altering the default listening-only workflow.
- [ ] Follow reference playback → ready to send → capture → finish → replay/reference/retry. Provide explicit keyboard controls to begin/end an attempt, retry, navigate and leave; do not guess completion solely from a pause in sending.
- [ ] Keep Iambic B timing and immediate sidetone on the adapter. Capture the resulting marks/gaps through the verified MIDI boundary; avoid duplicate live computer sidetone.
- [ ] Replay the measured mark/gap timing through existing audio output, using a chosen practice tone. Do not quantize or correct an attempt into ideal Morse before replaying it. Explain that this reconstructs timing rather than recording the adapter's sound.
- [ ] Keep the original reference timing, including Farnsworth spacing, independent of the adapter's sending speed. Do not silently synchronize or change hardware settings.
- [ ] Reuse pause/stop cleanup and generation checks so reference audio, replay and speech do not overlap. Gate capture by exercise state so touches during reference playback do not become an attempt.
- [ ] Preserve hidden/reveal behavior, reset attempt state on item/line changes, and keep automatic list/line progression from skipping a learner's turn.
- [ ] Handle empty attempts, long attempts, incomplete marks, cancellation and disconnect with bounded capture and an explicit incomplete state; never present a damaged capture as a valid attempt.
- [ ] Use recorded MIDI events to verify replay timing and interruptions, then demonstrate listen → touch-pad send → replay on the actual adapter.

## Boundary

This slice delivers capture and comparison by ear. Decoding and timing feedback
belong to 15. Attempts can remain session-local; persistent recordings, scoring,
radio output and microphone decoding are not implied by this ticket.
