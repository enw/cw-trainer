# 09 — Provide a curated editable starter library

**What to build:** Begin practice immediately with curated editable lists for KE2ISH, CW terms, prosigns, and realistic exchanges.

**Blocked by:** 08 — Edit exchanges and advance automatically through lines.

**Status:** implemented — curated seed and preservation/upgrade tests passed

- [x] Seed separate editable lists for KE2ISH, common terms, common prosigns, and short exchanges on first use; do not reseed over user edits or deliberately emptied lists.
- [x] Include CQ, DE, K, R, TU, 73, QTH, RST, and QRS, with useful additional terms and correct editable spoken phrases.
- [x] Include prosign meanings and verified continuous encodings; keep ordinary abbreviations distinct.
- [x] Include short ordered exchanges for calling CQ, answering, reports, name/location, repeats, slower sending, and signing off; use KE2ISH as the learner station and a consistently identified example counterpart.
- [x] Verify terminology, prosign use, and exchange wording against authoritative amateur-radio references; record sources and explain example callsigns are practice material.

## Boundary

This slice owns content and non-destructive first-run seeding. Ticket 10 owns the final integrated walkthrough and user-facing operating instructions; every feature ticket retains its own correctness checks.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. Sources and practice-example qualifications are in docs/design/content-sources.md. Existing installations can explicitly import fresh starter copies without modifying their old lists.
