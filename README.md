# cw-looper

A macOS terminal app for recognizing familiar CW sounds: your callsign,
common terms, continuous prosigns, and two-station exchanges. Start with
**KE2ISH**, repeat one item until ready, and optionally reinforce it with speech.

## Run

Requires Go 1.26+, macOS, Xcode Command Line Tools (for the native MIDI bridge),
`/usr/bin/say`, and a configured audio output.
Tested on macOS 26.6.2, Apple silicon. Use a terminal at least 80 × 24.

```sh
go run ./cmd/cw-looper
# Or build and run:
go build -o bin/cw-looper ./cmd/cw-looper
./bin/cw-looper
```

First launch installs four editable lists: callsign, common terms, prosigns,
and exchanges. Enter opens a list; Enter on an ordinary item starts practice.
For an exchange, open it and then start a selected line. Browsing is silent.
An eight-line example contact is included. Names, reports and locations are
practice examples, not a real contact; [content sources](docs/design/content-sources.md).

## Controls and editing

| Control | Browse | Practice |
| --- | --- | --- |
| ↑/↓ or k/j | Select | — |
| Enter or l | Open/start | Enter retries a failure |
| ←/→ or h/l | Back/open | Previous/next item or exchange line |
| [ / ] | — | Previous/next top-level item/exchange |
| Space | — | Pause/resume |
| s / S | — | Speak current phrase now |
| r / R | — | Toggle answer reveal |
| m | MIDI device diagnostic | — |
| t | Send back selected item/line (requires MIDI) | — |
| p | List settings (inside a list) | Settings overlay; playback continues |
| a | Add list, ordinary item or exchange line | — |
| A | Add an exchange inside a list | — |
| e / x | Edit / delete selected material | — |
| J / K | Move selected material down / up | — |
| i | Import fresh starter copies from Lists | — |
| Esc | Back | Stop and return to library |
| ? | Help | Help; playback continues |
| q | Quit | Stop and quit |
| Ctrl+C | Emergency exit | Stop and exit; report unsaved changes |

Forms use Tab / Shift+Tab to move fields, Ctrl+S to accept/save, and Esc to
cancel. Ordinary shortcut letters are text while a field has focus; pasting
works too. Ctrl+P in a CW field opens the prosign picker. Confirm deletion with
y; Esc/n cancels. To edit exchange lines, open the exchange first; e on its
list entry renames it. An exchange must retain at least one line.

Spoken overrides are optional. Clear the spoken field to use generated speech:
CQ is “see cue,” callsigns are spelled, words are read naturally, and prosigns
use their meanings. Custom speech survives CW edits, so review it after changing
CW. Overrides and all accepted settings save automatically.

## Practice settings

Type the listed choice in its settings field, then Ctrl+S to save the form.
Tab scrolls to the remaining fields in shorter terminals.

| Setting | Default | Choices / bounds |
| --- | --- | --- |
| Character / overall speed | 20 / 15 WPM | 5–60, overall ≤ character |
| Your / other tone | 440 / 600 Hz | 100–2000 Hz |
| Pause | 2 seconds | 0–30 seconds |
| Display | `whole` | `whole`, `highlight`, `hidden` |
| Playback | `cw` | `cw`, `sandwich`, `occasional` |
| Speech interval | 5 CW repetitions | 1–100 |
| Sequence | `repeat` | `repeat`, `cycle`, `shuffle` |
| Exchange lines | `manual` | `manual`, `automatic` |

Farnsworth timing keeps characters at character speed and stretches the
spacing to achieve overall speed. Equal speeds use standard spacing.

- `cw`: CW only. `sandwich`: CW → speech → CW. `occasional`: N complete CW
  repetitions → speech, beginning with CW. One configured gap separates portions.
- `repeat` stays on the current item. `cycle` wraps. `shuffle` visits each
  top-level item once per round. Advancement waits for the entire practice unit.
- Manual exchange lines repeat until Next, even with cycle/shuffle. Automatic
  lines retain conversational order; after the last line, repeat restarts that
  exchange, while cycle/shuffle selects the next top-level entry.
- `highlight` brackets the current character or entire prosign using estimated
  PCM consumption. `hidden` conceals the answer and title until R; reveal lasts
  across repetitions but resets on item/line changes.
- Pause restarts interrupted CW or speech from that portion's beginning. Gaps
  retain their remaining duration. Navigation stays silent while paused.
- Speak-now cancels the current portion, resets its unit, speaks, then resumes
  CW after the gap. From paused, it speaks once and returns to paused.
- Speed/tone changes apply at the next CW rendering; interval changes at the
  next practice unit. Style/sequence/line-mode changes restart the unit through
  cancellation, preserving pause. Display changes apply immediately.

## Vail touch-pad send-back practice

Connect your adapter, press **m** while browsing, select its MIDI input, then
**Enter** to enable MIDI. This changes only USB output mode; the adapter keeps
its keyer, speed and sidetone settings. The diagnostic shows received key events
and mark durations. **d** disconnects; **f** refreshes devices; **Esc** returns
with the connection active. Normal app exit returns the adapter to keyboard
output. Device removal requires reconnecting explicitly; a lost connection
invalidates capture instead of guessing missing marks.

Select an ordinary item or exchange line and press **t**:

1. Hear the reference once. Automatic list progression is suspended.
2. Release both pads, then **Enter** to begin capture. Send with local sidetone.
3. **Enter** finishes; **y** replays your measured timing, **v** the reference.
4. **Enter** starts another attempt; **n** starts a new reference and attempt.

Space pauses audio; pausing capture abandons the partial attempt and resumes at
Ready. h/l and arrows navigate; R reveals hidden text; Esc returns to browsing.
Two-minute/1024-mark limits bound capture. Empty, overlapping, unmatched or held
marks at finish are incomplete and cannot replay. Attempts are session-local.
The app captures MIDI timings, not microphone audio, and generates no duplicate
live sidetone. Hardware sending speed remains independent of reference WPM.
The [hardware record](docs/design/vail-midi-verification.md) describes the tested
QT Py M0 adapter's note-1/2 behavior; note-0 timed firmware is supported too.
Decoding and graded feedback are a later slice.

## Saved material and recovery

Data lives outside the binary, so replacing the application preserves it:

```text
~/Library/Application Support/cw-looper/library.json
```

Use `--data /path/to/library.json` for a separate library. Existing skeleton
installations remain unchanged; press i in Lists to explicitly add fresh starter
copies. Imports never replace existing lists. Deliberately empty libraries/lists
remain empty. Only one app may open a library at a time.

Supported CW: A–Z, 0–9, spaces, and `<AR>`, `<AS>`, `<BT>`, `<SK>`, `<HH>`.
Lowercase is normalized. Prosigns are continuous; ordinary `AR` has a letter
gap. Punctuation and unknown tokens are rejected. Limits: 512 bytes of CW,
1024 bytes of spoken text, 100 lists, 100 entries per list, 100 lines per exchange.
An ordinary item's title is a shortened preview; its CW text is preserved in full.

Invalid/future-version files are reported and preserved. Correct the file or
restore your backup, then relaunch. Saves use temporary-file sync and atomic
rename. A save failure leaves edits in memory marked **Not saved**; Ctrl+S retries.
Normal quit waits for a pending save and confirms discarding known unsaved edits;
Ctrl+C exits without trapping you behind confirmation and reports unsaved changes.
A companion `.lock` file can remain after exit; the OS releases its lock. Do not
remove it to bypass another running app.

Audio errors may require selecting a working system output and restarting.
For a speech failure, retry with Enter or set playback to `cw`. A failing backend
stops instead of looping. Speech text goes directly to macOS `say` on stdin;
no shell evaluates it. Audio/speech need no network service. Switching audio
devices mid-session is not verified; restart afterward.

## Verification and remaining acceptance

```sh
go test ./...
go test -race ./...
go vet ./...
# Deliberately plays tones and speech on the Mac:
CW_REAL_AUDIO=1 go test -v ./internal/output -run TestMacOutputLifecycle -count=1
```

The [complete verification record](docs/design/complete-verification.md) separates
automated checks, actual Mac output and PTY observations from human listening.
Audio quality, perceived highlighting synchronization, speech pronunciation and
acoustic stop tails still need listening acceptance. API cleanup timing cannot
prove silence at the speaker.

The [dependency map](.scratch/cw-practice/map.md) indexes the completed 00–10
implementation and queued 11–12 follow-ups. See also the [UI design](docs/design/terminal-ui.md)
and [architecture](docs/design/playback-architecture.md).
