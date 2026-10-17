# 10 — Verify and document a complete Mac practice session

**What to build:** Run the assembled application through a complete first-use and returning-user session, with reproducible instructions and honest evidence of what was verified.

**Blocked by:** 09 — Provide a curated editable starter library.

**Status:** implemented — software, real Mac lifecycle and PTY checks passed; human listening pending

- [x] Consolidate actual install/run instructions, data location, editing, settings, keyboard controls, supported input, and recovery from reported audio/speech/save failures.
- [x] Run the automated suite and a real Mac walkthrough covering first run, saved edits, restart, station tones, arrow/vi controls, hidden/reveal and highlighting, reinforcement patterns, cycle/shuffle, and automatic exchange advancement.
- [x] Check setting changes and interruptions during CW/speech/gaps, keyboard focus during editing, answer concealment, and terminal restoration on exit.
- [x] Demonstrate installation replacement or upgrade preserves saved material and settings, including intentionally empty lists.
- [x] Record automated evidence separately from real audio/listening/visual checks. Explicitly identify any checks requiring user observation instead of claiming unperformed validation.
- [x] Resolve integration defects in their owning feature implementation; update operating instructions to match the behavior actually delivered.

## Boundary

This is integrated acceptance and documentation, not a place to defer feature correctness, timing validation, or basic real-device checks.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. See docs/design/complete-verification.md for observed checks and remaining human listening acceptance.

## Remaining human acceptance

- [ ] Listen for tone quality/clicks, speech pronunciation, station contrast, perceived highlighting alignment and acoustic stop tails. Software/PTY checks cannot establish these.
