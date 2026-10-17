# 05 — Follow along or listen without seeing the answer

**What to build:** Choose whole-item display, synchronized character highlighting, or hidden text with manual reveal, independently of playback style.

**Blocked by:** 03 — Tune and remember each list's practice settings; 04 — Practice continuous prosigns and mixed phrases.

**Status:** draft — design and skeleton must inform feature readiness

- [ ] Default to whole-item display; let the learner select and persist another display mode for each list.
- [ ] Highlight the current ordinary character using the audio timing basis. Highlight an entire prosign as one unit.
- [ ] Hidden mode conceals answer text in the practice view, including titles/previews/progress; R reveals the current item.
- [ ] Design decision: a revealed answer remains visible across repetitions and hides when the item changes, including navigation back to it.
- [ ] Keep highlighting correct across pauses, cancellations, and next-repetition speed changes.
- [ ] Verify visible output via session events and a real terminal synchronization check; document limits of automated timing evidence.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is authored but awaits user review and skeleton validation; this ticket
has not started.
