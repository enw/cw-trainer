# Terminal interaction design

Ticket 00 design, 2026-10-16. Ready for user review; not implemented.
This resolves the spec's open interaction defaults. Product scope is unchanged.
Technical details live in [playback architecture](playback-architecture.md);
validation and the ticket 01 fixture live in [skeleton acceptance](skeleton-acceptance.md).

## Layout and focus

Use one focused view at a time. Enter opens/starts; Esc goes back and stops
practice. A practice screen never contains a library preview. Use text labels
and a visible focus marker, not color alone. Target 80 × 24 initially; below
that size stop playback and show a resize message, preserving selection.
Do not start playback automatically when launching or returning from an editor.
The sketches show the completed product; ticket 01 implements only the subset
specified in skeleton acceptance.

### Library: lists, then items

```text
cw-looper                                      Saved
Lists
> My callsign                 1 item
  Common CW                  18 items
  Prosigns                    8 items
  Exchanges                   4 exchanges

j/k or arrows: move   Enter/l: open   a: add   e: rename
x: delete   ?: help   q: quit
```

Counts are illustrative, not a content commitment. Opening a list replaces the
list view with its items; Esc/h returns. Entries show kind (item/exchange).
Enter starts an ordinary item or opens an exchange's ordered line selector;
Enter on a line starts practice there. j/k and arrows move selection. Library
selection itself is silent. Empty states explain the next action: “No items.
Press a to add one.” Delete uses an inline confirmation, with cancel selected
initially. J/K move the selected entry up/down in editing lists.

### Practice

```text
cw-looper                 Exchange 1 of 4 · Line 2 of 3
Station: You                                  440 Hz

              CQ CQ DE KE2ISH K

CW · Repetition 3                    Repeat current
20 character / 15 overall WPM         Pause 2.0 s
Display: Whole item                   Speech: CW only

Space pause   h/l previous/next line   [/]: previous/next item
r reveal   s speak   p settings   Esc library   q quit
```

In an ordinary item, h/l means previous/next item. In an exchange it means
previous/next line, clamped at the first/last line; show “Last line — ] next
item or Esc library.” The explicit [/] commands move between top-level items
from either kind, opening an exchange at line one. These manual commands clamp
at list ends; cycle wraps only during automatic playback. j/k do nothing on
the practice view, avoiding accidental navigation while retaining vi keys in
selection lists. Arrow left/right mirror h/l, up/down mirror j/k.

Whole-item display shows the complete current line. Highlight mode emphasizes
the current character and its containing word; a prosign is one highlight
unit. This preserves a whole-sound focus while making word/letter progress
available. Display progress is about CW, not speech word alignment.

Hidden mode replaces content with “Listen… r to reveal” and uses neutral
list/item/exchange labels plus numeric positions. Hide titles, callsigns,
previews, token counts/lengths, and character progress; role labels You/Other
and frequencies may remain. Explicit speech can reveal meaning by ear. r
reveals for this item/line until it changes; selecting it again hides it.
Changing into hidden mode immediately hides it. Help/settings overlays must
not expose hidden content. Leaving practice stops it before the library is
shown; browsing the library intentionally exposes its contents.

### Item and exchange editors

```text
Edit item                                      Not saved
CW text       [CQ CQ DE KE2ISH <AR>                    ]
Spoken text   [see cue see cue D E K E 2 I S H ...       ]
              Generated default · Override available

Tab next field   Ctrl+P prosign picker   Ctrl+S save   Esc cancel
```

While a field has focus, ordinary keys—including h/j/k/l/r/s/q and Space—are
text. Arrows move the cursor. Tab/Shift+Tab change fields; Ctrl+S validates and
saves; Esc cancels the uncommitted form. This is automatic persistence of every
accepted edit, not saving partial invalid keystrokes. Reordering persists
immediately. Invalid input leaves the form open with its draft intact.

A prosign picker shows token and meaning; Enter inserts it at the cursor.
Generated speech follows CW edits only while no override exists; editing the
spoken field creates an override. An explicit “Use generated speech” action
clears it. Editing CW does not overwrite a custom phrase; show a reminder to
check it. Ticket 01 uses explicit fixture phrases without this editor.

An exchange editor first shows ordered lines with role and text; a/e/x/J/K
add/edit/delete/reorder. A line form adds a You/Other selector to the item form.
Save disallows blank text and an empty exchange. Pause/stop practice before
opening any content editor. Returning from it leaves playback stopped; Enter
starts with the current saved content. Deleting the selected entry selects its
successor, or predecessor if it was last, or the empty state.

### Settings

```text
Practice settings                              Saved
Character WPM        [20]    Overall WPM       [15]
Your tone Hz         [440]   Other tone Hz     [600]
Pause seconds        [2.0]   Speak every N     [5]
Playback             [CW only / CW-talk-CW / Occasional]
Display              [Whole item / Highlight / Hidden]
Sequence             [Repeat / Cycle / Shuffle]
Exchange lines       [Manual / Automatic]

Tab move   Enter apply field/choice   Esc practice
Playing: 20/15 · 440 Hz     Pending next CW: 22/15 · 500 Hz
```

Settings is an overlay; practice continues unless Space is pressed outside a
text field. Accepted values save immediately; incomplete numeric input stays
local. Speed/tone changes wait for the next complete CW rendering. Pause changes
affect the next newly scheduled gap. Style/sequence/line-mode changes cancel the
current unit and start a fresh one on the same item, preserving paused state.
Interval changes apply to the next unit, with the active interval visible until
then. Display changes apply immediately. Esc discards an unaccepted field only.

Initial bounds: character/overall speeds 5–60 WPM (overall ≤ character), tone
100–2000 Hz, pause 0–30 s, speech interval 1–100. These are initial product bounds,
not physiological claims; ticket 01 checks audio assumptions before ticket 03
implements them. New lists start CW-only, whole-item, repeat-current, manual
lines, 20/15 WPM, 440/600 Hz, 2 s pause, and interval five.

## Playback policies

| Situation | Chosen behavior |
| --- | --- |
| Practice unit | CW only: one CW; CW-talk-CW: CW, speech, CW; occasional: N CW, speech |
| Separation | One configured gap between portions and between units; no double gap at a unit boundary; no initial gap |
| Repeat | Repeat current item/line after each unit |
| Occasional speech | Start with CW; count completed renderings only; reset after speech, navigation, or style change |
| On-demand speech | Cancel current portion, reset the current unit, speak once, then gap and CW; it does not increment the occasional counter |
| Pause during CW/speech | Stop output; resume restarts that portion from its beginning; retain already completed portions/counts |
| Pause during a gap | Freeze remaining duration; resume continues that remainder |
| Speak while paused | Speak once, then remain paused with next resume starting a fresh unit |
| Navigate while paused | Select new target silently; next resume starts a fresh unit |
| Reveal | Persists over repetitions; resets on item/line change |
| Cycle | Advance top-level items in order and wrap |
| Shuffle | Each top-level item once per round; no special rule across round boundaries |
| Exchange manual | Repeat line until Next, even if top-level sequence is cycle/shuffle |
| Exchange automatic | Advance ordered lines after each unit; after last line apply top-level sequence |
| Automatic lines + repeat | Repeat the current exchange from line one after its final line |

Manual line mode thus gates top-level automatic advancement while inside an
exchange. Show this explicitly in help: “Manual lines: press Next.” Manual
selection/navigation resets the current unit and shuffle round, beginning the
new round with the selected top-level entry. Pause alone does not reset either.

## Errors and recovery

Playback failure stops the session, shows a neutral error without hidden text,
and offers retry or Esc to library. Never keep looping a failing backend.
A fatal audio-context failure requires app restart; say so rather than promising
an in-process retry. Speech failure offers retry or selection of CW-only mode.
Load corruption opens a read-only error view with the data path and recovery
instructions; never seed over the file. Save failure leaves the edited data in
memory marked “Not saved” with retry; quit confirms discarding unsaved changes.
Ctrl+C always stops output and exits promptly; persist accepted edits beforehand
and report unsaved data rather than trapping emergency exit behind a dialog.
