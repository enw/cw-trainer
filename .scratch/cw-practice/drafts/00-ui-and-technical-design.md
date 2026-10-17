# 00 — Design the terminal workflow and playback architecture

**What to build:** Produce a lightweight, reviewable UI and technical design that makes the walking skeleton implementable and resolves the remaining interaction choices.

**Blocked by:** None — first work item when design work is requested.

**Status:** draft — design and skeleton must inform feature readiness

- [ ] Sketch the library, practice, item/exchange editor, and settings views, including empty/error states, focus, keyboard help, and navigation while playing or editing.
- [ ] Define arrow/vi mappings and distinguish next item, next exchange line, and leaving an exchange. Show station roles and how hidden-answer mode avoids previews or title leaks.
- [ ] Resolve the spec’s proposed defaults: reveal lifetime, speech counts, practice-unit completion, shuffle boundaries, speak-now interruption, and pause/resume behavior. Record decisions in the design; flag any change to approved requirements for user review.
- [ ] Choose and justify a minimal language, terminal UI, audio, speech, and persistence approach using current primary documentation. Record prerequisites and the riskiest assumptions for the skeleton to test.
- [ ] Describe one owner for playback timing/transitions, with a state diagram covering CW, silence, speech, paused, cancellation, completion, and error recovery. Define how display progress follows audio and how setting changes are deferred.
- [ ] Define a saved model for lists, ordinary items, ordered exchange lines, station roles, spoken overrides, and per-list settings. Include data location, durable saves, first-run seeding, and preservation of existing user data.
- [ ] Define the practice-session verification boundary, deterministic clock/output controls, and the real Mac checks needed to validate audio, speech, and responsive keyboard input.
- [ ] Specify a small two-station skeleton fixture and its demo steps. Review the design with the user before implementing it; remaining stack/audio assumptions are explicitly provisional until ticket 01.
- [ ] Update this same draft set with resolved decisions and dependencies; do not create an alternative implementation plan.

## Boundary

This is design work, not application implementation. Keep the artifacts small: terminal wireframes, playback state diagram, data model, and a short stack decision. No separate throwaway spike is planned; ticket 01 validates the choices with working application code.
