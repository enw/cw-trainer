# 06 — Reinforce CW with editable spoken prompts

**What to build:** Hear CW → speech → CW, or speech after a configurable number of CW repetitions, and request speech immediately with S.

**Blocked by:** 03 — Tune and remember each list's practice settings; 04 — Practice continuous prosigns and mixed phrases.

**Status:** draft — breakdown awaiting review

- [ ] Offer CW only, CW → speech → CW, and occasional-speech modes; persist the selected mode and speech interval per list.
- [ ] Default occasional speech to every five completed CW repetitions, beginning with CW; validate a positive integer interval.
- [ ] Use built-in Mac speech. Read words naturally, CQ as “see cue,” callsigns character by character, and prosigns by meaning, combining these rules for full phrases.
- [ ] Allow an editable spoken override for every item; preserve it when loading and clearly handle subsequent CW-text edits.
- [ ] Keep the shared configurable pause between portions, starting after speech completes. Never overlap speech and CW.
- [ ] S speaks the current phrase on demand. Proposed default: interrupt the current portion, speak, then restart CW after the pause; when already paused, remain paused afterward.
- [ ] Proposed default: reset the repetition count on item/style changes; count only completed CW repetitions. Restart a paused spoken portion if native resumption is unavailable.
- [ ] Pause, navigation, and exit cancel or suspend speech cleanly; failures are visible and recoverable. Pass speech text as data, never shell syntax.
- [ ] Verify portion order, counts, cancellation, speech failures, and pause timing with controlled outputs; smoke-check actual Mac speech and pronunciation.

## Boundary

Speech does not require a particular display mode. Cross-mode coverage is completed once ticket 05 is also available.
