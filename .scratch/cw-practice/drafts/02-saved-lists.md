# 02 — Create and practice saved message lists

**What to build:** Create named lists and edit, reorder, select, and practice ordinary text items entirely in the terminal, with changes surviving restart.

**Blocked by:** 01 — Run a saved two-station walking skeleton.

**Status:** implemented — verification included in final integration

- [x] Create and rename lists; add, edit, remove, and reorder ordinary message items. Retain spoken-phrase data without adding the pronunciation editor owned by ticket 06.
- [x] Extend the skeleton’s durable storage to autosave list/item edits; surface save failures and preserve existing data on read/parse errors.
- [x] Navigate lists with arrows and j/k; move between practice items with left/right and h/l. Shortcuts do not consume typing in text fields.
- [x] Repeat the selected item until navigation; navigation cancels prior output without overlapping the newly selected item.
- [x] Handle empty lists, deletion of the active item, and first/last-item navigation without crashes or stale playback. Document boundary behavior.
- [x] Reject unsupported ordinary characters with a useful message instead of silently dropping them. Document the supported repertoire.
- [x] Verify edit/save/reload behavior, a meaningful failed-save case, and navigation during tone/gap output through user-level actions.

## Boundary

Reuse the skeleton’s loading, saving, playback, and cancellation behavior. Exchange line editing belongs to ticket 08. Only minimal initial material is needed here. The complete curated library is delivered in ticket 09.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. This ticket has not started.
