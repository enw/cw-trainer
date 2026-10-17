# CW practice dependency map

This indexes the single current draft set; the linked tickets own their scope
and acceptance criteria. The spec remains the product requirements source.
Design and skeleton findings should revise these tickets in place.

| Ticket | Delivers | Blocked by |
| --- | --- | --- |
| [00](drafts/00-ui-and-technical-design.md) | UI and technical design | None |
| [01](drafts/01-walking-skeleton.md) | Real audio/speech/storage through one two-station exchange | 00, including design review |
| [02](drafts/02-saved-lists.md) | List and ordinary-item editing | 01 |
| [03](drafts/03-adjust-practice-settings.md) | Runtime timing and both station tones, saved per list | 02 |
| [04](drafts/04-prosigns-and-phrases.md) | Prosign picker and continuous mixed-phrase playback | 02 |
| [05](drafts/05-display-and-reveal.md) | Whole-item, highlighting, and hidden/reveal display | 03, 04 |
| [06](drafts/06-spoken-reinforcement.md) | Spoken-phrase editing and reinforcement patterns | 03, 04 |
| [07](drafts/07-cycle-and-shuffle.md) | Cycle/shuffle complete practice units | 06 |
| [08](drafts/08-two-station-exchanges.md) | Exchange editing and automatic line advancement | 05, 07 |
| [09](drafts/09-starter-library-and-first-run.md) | Curated starter material and safe seeding | 08 |
| [10](drafts/10-mac-integration-and-usage.md) | Integrated Mac walkthrough and usage instructions | 09 |

00 resolves interactions and selects the stack. 01 validates those choices
with actual application code before the remaining feature work. There is no
separate disposable prototype or second implementation plan.

03/04 and later 05/06 are independent dependency branches, not an instruction
to delegate or implement concurrently. Each ticket verifies its own behavior;
10 checks the assembled experience rather than absorbing deferred correctness.

Ticket 00 design is approved and resolved. Ticket 01 is implemented and
ready-for-human for final listening acceptance; see the
[verification record](../../docs/design/skeleton-verification.md).
Tickets 02–10 have not started. The current user request authorized only the
walking skeleton. Software checks do not replace audible quality/tail validation.

## Queued follow-ups

The user requested these issues during implementation; they are not part of the
current 02–10 build. They are independent of each other and reuse completed 09.

| Ticket | Delivers | Blocked by |
| --- | --- | --- |
| [11](issues/11-nato-phonetic-alphabet.md) | NATO/ITU phonetic alphabet practice and reference | 09 |
| [12](issues/12-rst-system.md) | RST report practice and rating-scale reference | 09 |
