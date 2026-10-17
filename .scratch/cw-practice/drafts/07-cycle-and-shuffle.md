# 07 — Cycle or shuffle complete practice units

**What to build:** Switch between repeating the current item, cycling a list, and shuffling it while preserving each item's chosen reinforcement pattern.

**Blocked by:** 06 — Reinforce CW with editable spoken prompts.

**Status:** implemented — deterministic unit-boundary and shuffle tests passed

- [x] Keep repeat-current as the default; expose cycle and shuffle and save the selection per list.
- [x] Design decision: CW-only advances after one rendering; CW → speech → CW advances after its final CW; occasional-speech advances after its configured CW repetitions and spoken reinforcement.
- [x] Design decision: cycle wraps; shuffle visits each item once per round before starting another. No guarantee against repeats across round boundaries is required.
- [x] Pause and exit halt progression; manual navigation cancels the current unit and resets its repetition count.
- [x] Handle empty and single-item lists and list edits during practice without runaway advancement or stale item references.
- [x] Verify progression and counts with a controllable clock and deterministic shuffle input, including all playback styles and list sizes.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. Content editing is available only after leaving practice, so no active snapshot can be deleted.
