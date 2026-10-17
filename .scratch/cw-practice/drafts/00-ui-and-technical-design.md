# 00 — Design the terminal workflow and playback architecture

**What to build:** Produce a lightweight, reviewable UI and technical design that makes the walking skeleton implementable and resolves the remaining interaction choices.

**Blocked by:** None — first work item when design work is requested.

**Status:** ready-for-human — design authored; user review pending

- [x] Sketch the library, practice, item/exchange editor, and settings views, including empty/error states, focus, keyboard help, and navigation while playing or editing.
- [x] Define arrow/vi mappings and distinguish next item, next exchange line, and leaving an exchange. Show station roles and how hidden-answer mode avoids previews or title leaks.
- [x] Resolve the spec’s proposed defaults: reveal lifetime, speech counts, practice-unit completion, shuffle boundaries, speak-now interruption, and pause/resume behavior. Record decisions in the design; flag any change to approved requirements for user review.
- [x] Choose and justify a minimal language, terminal UI, audio, speech, and persistence approach using current primary documentation. Record prerequisites and the riskiest assumptions for the skeleton to test.
- [x] Describe one owner for playback timing/transitions, with a state diagram covering CW, silence, speech, paused, cancellation, completion, and error recovery. Define how display progress follows audio and how setting changes are deferred.
- [x] Define a saved model for lists, ordinary items, ordered exchange lines, station roles, spoken overrides, and per-list settings. Include data location, durable saves, first-run seeding, and preservation of existing user data.
- [x] Define the practice-session verification boundary, deterministic clock/output controls, and the real Mac checks needed to validate audio, speech, and responsive keyboard input.
- [x] Specify a small two-station skeleton fixture and its demo steps; keep stack/audio assumptions explicitly provisional until ticket 01.
- [ ] Review the completed design with the user before ticket 01 implementation.
- [x] Update this same draft set with resolved decisions and dependencies; do not create an alternative implementation plan.

## Boundary

This is design work, not application implementation. Keep the artifacts small: terminal wireframes, playback state diagram, data model, and a short stack decision. No separate throwaway spike is planned; ticket 01 validates the choices with working application code.

## Deliverables

- [Terminal UI](../../../docs/design/terminal-ui.md): wireframes, keys, defaults, focus, error states.
- [Playback architecture](../../../docs/design/playback-architecture.md): sourced stack choice, state diagram, model, durability, audio risks.
- [Skeleton acceptance](../../../docs/design/skeleton-acceptance.md): fixture, demo, agreed verification seams.

No application code or dependency installation was performed. User review remains pending; ticket 01 has not started.

## Verification

- Local Markdown links, fenced blocks, file-size limits, dependency ordering,
  and Git whitespace checks passed.
- Independent standards review of `259f217..9db06a6`: zero findings.
- Independent spec/scope review of the same change: zero material findings.
- No application tests, rendered-terminal review, or audio experiments were run;
  these remain walking-skeleton checks. Agent reviews do not replace user design review.
