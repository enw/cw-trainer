# 06 — Reinforce CW with editable spoken prompts

**What to build:** Hear CW → speech → CW, or speech after a configurable number of CW repetitions, building on the skeleton’s working S speak-now control.

**Blocked by:** 03 — Tune and remember each list's practice settings; 04 — Practice continuous prosigns and mixed phrases.

**Status:** implemented — controlled portion-order and cancellation tests passed

- [x] Offer CW only, CW → speech → CW, and occasional-speech modes; persist the selected mode and speech interval per list.
- [x] Default occasional speech to every five completed CW repetitions, beginning with CW; validate a positive integer interval.
- [x] Reuse the skeleton’s built-in Mac speech integration. Generate default phrases that read words naturally, CQ as “see cue,” callsigns character by character, and prosigns by meaning, combining these rules for full phrases.
- [x] Provide a spoken-phrase editor for ordinary items and existing exchange lines; allow an editable override for each; preserve it when loading and clearly handle subsequent CW-text edits.
- [x] Keep the shared configurable pause between portions, starting after speech completes. Never overlap speech and CW.
- [x] Preserve the skeleton’s reviewed speak-now/pause/navigation semantics across the new playback modes.
- [x] Design decision: reset the repetition count on item/style changes; count only completed CW repetitions. Follow the skeleton’s reviewed speech-resumption behavior.
- [x] Extend the skeleton’s cancellation and error handling across every new mode transition; verify speech and CW never overlap.
- [x] Verify portion order, counts, cancellation, speech failures, and pause timing with controlled outputs; actual Mac speech is exercised in ticket 10; human pronunciation acceptance remains pending.

## Boundary

Speech does not require a particular display mode. Cross-mode coverage is completed once ticket 05 is also available.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. Software verification passed. Actual pronunciation and audio tails remain a listening acceptance item.
