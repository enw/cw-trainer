# cw-looper

## Starting brief

A command-line app for creating a list of Morse code / CW messages and playing
them as dits and dahs while showing the current word and/or letter.

- Configurable audio tone, starting at 440 Hz.
- Configurable sending speed, starting at 15 WPM.
- Language, terminal interaction, storage, playback controls, and precise timing
  behavior remain open for requirements discussion.

## Workflow

The user will request each stage: grilling, then spec, then tickets, then
implementation. Do not advance to the next stage without that request.
The ticket sequence starts with UI/technical design, followed by a working
two-station skeleton, then feature slices. The canonical dependency index is
`.scratch/cw-practice/map.md`; revise the existing drafts rather than keeping
alternative plans. Validate the design through the skeleton before treating
the remaining feature drafts as implementation-ready.
Commit completed work incrementally. Review status and stage only intended
files before each commit; do not push unless requested.
During grilling, ask one decision at a time. Capture approved requirements in
the spec before creating tickets or implementing the app.

Matt Pocock's skills are installed globally under `~/.agents/skills/`.
Read the relevant `SKILL.md` when invoking a stage:

- `grilling`: requirements discussion.
- `to-spec`: synthesize the agreed requirements.
- `to-tickets`: break the approved spec into implementation tickets.
- `implement`: implement a ticket.

No application code or build tooling exists yet. Record actual run, test, and
lint commands here once the implementation establishes them.

## Agent skills

### Issue tracker

Use local Markdown in `.scratch/<feature-slug>/`, with `spec.md` and one numbered
file per ticket under `issues/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Use the default five triage labels. See `docs/agents/triage-labels.md`.

### Domain docs

Use a single-context layout: root `CONTEXT.md` and `docs/adr/`, created when needed.
See `docs/agents/domain.md`.
