# 02 — Create and practice saved message lists

**What to build:** Create named lists and edit, reorder, select, and practice ordinary text items entirely in the terminal, with changes surviving restart.

**Blocked by:** 01 — Hear your callsign in a terminal practice loop.

**Status:** draft — breakdown awaiting review

- [ ] Create and rename lists; add, edit, remove, and reorder ordinary message items. Each item can retain an optional editable spoken phrase for later speech playback.
- [ ] Autosave changes durably outside the installed application; surface save failures and preserve existing data on read/parse errors.
- [ ] Navigate lists with arrows and j/k; move between practice items with left/right and h/l. Shortcuts do not consume typing in text fields.
- [ ] Repeat the selected item until navigation; navigation cancels prior output without overlapping the newly selected item.
- [ ] Handle empty lists, deletion of the active item, and first/last-item navigation without crashes or stale playback. Document boundary behavior.
- [ ] Reject unsupported ordinary characters with a useful message instead of silently dropping them. Document the supported repertoire.
- [ ] Verify edit/save/reload behavior, a meaningful failed-save case, and navigation during tone/gap output through user-level actions.

## Boundary

Only minimal initial material is needed here. The complete curated library is delivered in ticket 09.
