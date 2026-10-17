# 05 — Follow along or listen without seeing the answer

**What to build:** Choose whole-item display, synchronized character highlighting, or hidden text with manual reveal, independently of playback style.

**Blocked by:** 03 — Tune and remember each list's practice settings; 04 — Practice continuous prosigns and mixed phrases.

**Status:** implemented — PCM consumption progress, generation checks and hidden/reveal

- [x] Default to whole-item display; let the learner select and persist another display mode for each list.
- [x] Highlight the current ordinary character using the audio timing basis. Highlight an entire prosign as one unit.
- [x] Hidden mode conceals answer text in the practice view, including titles/previews/progress; R reveals the current item.
- [x] Design decision: a revealed answer remains visible across repetitions and hides when the item changes, including navigation back to it.
- [x] Keep highlighting correct across pauses, cancellations, and next-repetition speed changes.
- [ ] Human audio/visual synchronization acceptance; software uses generation-tagged PCM consumption progress.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. Display uses consumed PCM minus queued buffers, not a separate word timer. Device latency and perceptual synchronization still need listening acceptance.
