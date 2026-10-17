# Design verification and walking-skeleton handoff

Ticket 00 design, 2026-10-16. These are planned checks, not results from a running
application. [UI](terminal-ui.md) and [architecture](playback-architecture.md) are
the canonical designs. Ticket 01 remains blocked on user design review and an
explicit request to implement it.

## Small fixture

One saved list, one exchange, two lines; whole-line display, fixed default
settings, manual advancement, and repeat-current. No editors, prosigns, shuffle,
reinforcement modes, or character highlighting are required yet.

| Line | Role / tone | CW | Explicit spoken phrase |
| --- | --- | --- | --- |
| 1 | Other / 600 Hz | CQ CQ DE W1AW K | see cue see cue D E W one A W K |
| 2 | Learner / 440 Hz | W1AW DE KE2ISH K | W one A W D E K E two I S H K |

W1AW here is example text only, not a recording or claim of a real contact.
This minimal call/answer fixture tests station roles and familiar sounds; the
full procedural/content review belongs to ticket 09.

## Demonstration script for ticket 01

1. Build/install using the selected stack and record exact versions and commands.
2. Start without existing data. Confirm the fixture was saved and launch line 1
   explicitly; listen to its complete repeated sound at 600 Hz, 20/15 WPM.
3. Press Space during CW. Confirm output stops; resume restarts that rendering.
   Repeat while in a gap and check the remaining gap resumes.
4. Press s during CW. Confirm CW stops before speech starts, then a two-second
   pause precedes restarted CW. Repeat s while paused; it must return to paused.
5. Advance with right/l during CW and during speech. Confirm only line 2 now
   plays, at 440 Hz. Rapid alternating navigation must settle on the last target.
6. Pause and navigate left/h. Confirm line 1 is selected silently; resume starts
   it. Check the first/last line clamps and Esc stops before showing the library.
7. Quit during CW and during speech. Confirm no process/audio continues and
   the terminal accepts normal shell input afterward.
8. Without adding an editor, alter a fixture value through the storage test
   seam or by editing the closed app's JSON, reopen, and confirm it is retained.
   Exercise failed save and invalid existing data without destructive reseeding.
9. Record audible tails and perceived responsiveness on built-in or wired
   output. Target stop/drain within 100 ms; distinguish observed estimates from
   instrumented measurements. If not measurable, mark the check unverified.
10. Reconcile findings into the design and existing ticket map before feature work.

## Automated checks at agreed seams

- Drive the practice session with user actions and output completion/stop/error
  events. Observe snapshots/effects using controlled time and tagged output IDs.
  Confirm navigation, cancelled/stale completion rejection, pause/resume,
  speak-now from paused/running, and pending timing changes when introduced.
- Independently assert known Morse timing: E is one dit, T three dits; ET has
  an inter-character gap; PARIS with its trailing word gap totals four seconds
  at 20/15 and three seconds at 20/20. Check cumulative sample rounding.
- Inspect actual rendered PCM for duration/frequency/envelope behavior; verify
  the output adapter's buffer-tail policy separately from ideal session events.
- Save/reload through the store's public interface, including existing empty
  data, malformed/schema-too-new files, write failures, and competing instances.
- Integrate with real `say` to establish lifecycle/cancellation behavior when
  ticket 01 runs. Keep speech text out of command interpretation and hidden UI.

Use independent expected values and observable behavior rather than mirroring
the implementation. No application test suite is created for ticket 00 because
there is no application yet. File/link and dependency validation is sufficient
for this documentation change; it is not visual or audio QA.

## Ownership after the skeleton

- 02 adds list/item editing to the existing model/store.
- 03 adds adjustable timing/tones and persistence of those settings.
- 04 adds prosign validation, encoding, and picker.
- 05 adds display modes, word/letter highlighting, and hidden-answer checks.
- 06 adds generated/editable speech and reinforcement patterns to existing speech.
- 07 adds cycle/shuffle of complete units.
- 08 adds exchange editing and automatic line advancement to existing exchanges.
- 09 adds reviewed starter material and safe library installation.
- 10 performs integrated verification and usage documentation.

No dependency changes were needed after this design. The remaining drafts link
to these decisions, so future work has one source for interaction semantics.
