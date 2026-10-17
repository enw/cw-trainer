# 06 — Reinforce CW with editable spoken prompts

**What to build:** Hear CW → speech → CW, or speech after a configurable number of CW repetitions, building on the skeleton’s working S speak-now control.

**Blocked by:** 03 — Tune and remember each list's practice settings; 04 — Practice continuous prosigns and mixed phrases.

**Status:** draft — design and skeleton must inform feature readiness

- [ ] Offer CW only, CW → speech → CW, and occasional-speech modes; persist the selected mode and speech interval per list.
- [ ] Default occasional speech to every five completed CW repetitions, beginning with CW; validate a positive integer interval.
- [ ] Reuse the skeleton’s built-in Mac speech integration. Generate default phrases that read words naturally, CQ as “see cue,” callsigns character by character, and prosigns by meaning, combining these rules for full phrases.
- [ ] Provide a spoken-phrase editor for ordinary items and existing exchange lines; allow an editable override for each; preserve it when loading and clearly handle subsequent CW-text edits.
- [ ] Keep the shared configurable pause between portions, starting after speech completes. Never overlap speech and CW.
- [ ] Preserve the skeleton’s reviewed speak-now/pause/navigation semantics across the new playback modes.
- [ ] Proposed default: reset the repetition count on item/style changes; count only completed CW repetitions. Follow the skeleton’s reviewed speech-resumption behavior.
- [ ] Extend the skeleton’s cancellation and error handling across every new mode transition; verify speech and CW never overlap.
- [ ] Verify portion order, counts, cancellation, speech failures, and pause timing with controlled outputs; smoke-check actual Mac speech and pronunciation.

## Boundary

Speech does not require a particular display mode. Cross-mode coverage is completed once ticket 05 is also available.
