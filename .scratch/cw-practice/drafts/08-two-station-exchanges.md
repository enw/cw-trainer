# 08 — Edit exchanges and advance automatically through lines

**What to build:** Extend the skeleton’s two-station line practice with full exchange editing and optional automatic line advancement, integrating the completed practice modes.

**Blocked by:** 05 — Follow along or listen without seeing the answer; 07 — Cycle or shuffle complete practice units.

**Status:** implemented — exchange round-trip and ordered advancement tests passed

- [x] Create/edit/reorder/remove exchanges and their lines, assigning CW text and learner/other station roles; reuse the spoken-phrase editor from ticket 06 and preserve edits on restart.
- [x] Reuse the skeleton’s line selection and station routing, with both adjustable tones supplied by ticket 03. Do not create a second exchange model or playback controller.
- [x] Keep repeat-until-Next as default; add automatic line advancement after the current line’s complete practice unit, with behavior at exchange boundaries matching the reviewed design.
- [x] Preserve conversational line order when top-level items/exchanges are shuffled and support lists mixing ordinary items and exchanges.
- [x] Apply ticket 05 display/reveal and ticket 06 reinforcement behavior to the current line; reset reveal/count state on line changes without leaking answers or retaining cancelled output.
- [x] Handle editing/deleting the active line, empty exchanges, and single-line exchanges without stale output or invalid navigation.
- [x] Verify save/reload, ordered automatic/manual progression, station tones, and combinations of display, speech, and sequencing modes through the existing session boundary.

## Boundary

Basic exchange playback, station routing, persistence, and navigation are owned by the skeleton. This slice adds editing and automatic progression; it does not introduce exchanges for the first time.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. Content editors open only outside playback. Empty exchanges are rejected; deleting a final line directs the learner to delete its parent entry.
