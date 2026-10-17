<div align="center">

# cw-looper

### Hear the shape. Send it back.

**A little terminal practice. A lot of dits and dahs.**

![Platform: macOS](https://img.shields.io/badge/platform-macOS-111827?style=for-the-badge&logo=apple&logoColor=white)
![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Input: USB MIDI](https://img.shields.io/badge/input-USB_MIDI-8B5CF6?style=for-the-badge)
![Controls: vi keys](https://img.shields.io/badge/controls-vi_keys-22C55E?style=for-the-badge)

[Get started](#get-started) · [Test your Vail adapter](#test-your-vail-adapter) · [Controls](#controls-and-editing) · [Settings](#practice-settings) · [Development](#development-and-verification)

</div>

Learn to recognize **whole CW sounds**: your callsign, familiar words, continuous
prosigns, and full exchanges. Start with **KE2ISH**, stay on one item until it
feels familiar, then move on at your own pace.

| Listen | Send | Understand |
| :--- | :--- | :--- |
| Repeat familiar sounds with Farnsworth timing | Send back on Vail capacitive pads | Replay your actual timing |
| CW → speech → CW, or occasional spoken reminders | Keep Iambic B and the adapter's local sidetone | Inspect tentative decoding and spacing |
| Two-station exchanges with distinct tones | Practice one exchange line at a time | Reveal hidden answers when ready |

```text
LISTEN  ──→  READY  ──→  SEND  ──→  REPLAY + FEEDBACK
             Enter      Enter       y / f
```

> **Status:** Listening practice, Vail capture/replay, and tentative feedback are
> implemented. Automated checks pass; the live `A` capture/decode/replay walkthrough
> passed on 2026-10-17, as did reconnecting by reselecting QT Py M0. Broader live
> feedback checks remain open. Badges describe the stack, not CI status.

## Get started

Requires Go 1.26+, macOS, Xcode Command Line Tools (for the native MIDI bridge),
`/usr/bin/say`, and a configured audio output.
Tested on macOS 26.6.2, Apple silicon. Use a terminal at least 80 × 24.

From the repository folder:

```sh
go build -o bin/cw-looper ./cmd/cw-looper
./bin/cw-looper
```

For a quick run without a separate build step: `go run ./cmd/cw-looper`.

First launch installs six editable lists: callsign, common terms, prosigns,
exchanges, the NATO/ITU alphabet, and RST reports. Enter opens a list; Enter on an ordinary item starts practice.
For an exchange, open it and then start a selected line. Browsing is silent.
An eight-line example contact is included. Names, reports and locations are
practice examples, not a real contact; [content sources](docs/design/content-sources.md).

## Test your Vail adapter

**Start with a single `A` — dit-dah — so the result is easy to check.** Allow about
five minutes. Use the freshly built binary above; quit any older running copy first.

### 1. Prove the connection

1. Plug in the Vail adapter. Leave it in your usual **Iambic B** configuration.
2. While browsing, press **m**, select **QT Py M0** (or your adapter's name),
   then **Enter** to enable MIDI output.
3. Tap each pad. The diagnostic should show **down/up events**, an increasing
   event count, and **zero held notes after release**. At the previously observed
   15 WPM, dits should be about **80 ms**, dahs about **240 ms**.
4. Press **Esc** to return to the library with MIDI connected.

Hearing the adapter's buzzer alone is not this test: the event display confirms
that the Mac receives your keying.

### 2. Make a tiny practice target

1. Open any list with **Enter**. Press lowercase **a** to add an ordinary item.
2. In **CW text**, type `A`. Leave the spoken override empty and press **Ctrl+S**.
3. Select the new `A` item with **j/k** or arrows, then press **t** for send-back.
   Wait for the reference to finish and the screen to say **Ready to send**.

If you only see **First Exchange**, press **i** on the top-level Lists screen to
import the starter library. Existing saved material is preserved.

### 3. Send, finish, listen

1. Release both pads, then press **Enter** to arm capture.
2. Send **dit-dah** as one character. Keep the gap between elements short;
   at 15 WPM it is roughly 80 ms. Release the pads.
3. Press **Enter** to finish. Expect **2 captured marks** and a tentative decode
   of **A**. With the answer visible, expect **Matching Morse patterns and word
   boundaries**.
4. Press **y**. The computer should replay your dit-dah with the timing you sent.
5. After the attempt or its replay, press **f** for full feedback. **j/k** scrolls;
   **Esc** returns. **- / +** adjusts the interpretation's sending WPM. Match it
   to the adapter; this control does not change the adapter or reference speed.
6. Back on the practice screen, press **v** to hear the reference again. It may
   differ from your replay: reference Farnsworth speed and adapter sending speed
   are independent. **Enter** at Ready starts a new attempt.

> **What counts as a pass?** Events reach the Mac, the attempt contains both
> marks, the decode is `A` at the correct sending speed, and replay sounds like
> your sending. An uncertain decode is useful evidence too: note the sending WPM
> and gap values rather than adjusting them just to force a match.

### 4. Try the edge cases

| Try this | Expected result |
| :--- | :--- |
| Capture `T` against the `A` target | Decode `T`; different Morse patterns |
| Start and finish without touching a pad | No input; no replayable attempt |
| Set display to `hidden` with **p**, then capture | Your sending is shown; target/comparison waits for **R** |
| Unplug during capture | Connection lost; attempt invalidated |
| Reconnect, **Esc** to browsing, **m**, **f**, select device, **Enter** | Fresh connection; start a new attempt |
| Quit normally with **q** | App closes its MIDI ports and requests keyboard-output mode |

**Tell me what you see:** adapter name, captured mark count, decoded text,
sending WPM, whether replay sounds right, and any disconnect/reconnect problem.
That completes the human check that automated tests cannot provide.

<details>
<summary><strong>Something didn't work?</strong></summary>

| Symptom | Check |
| :--- | :--- |
| No MIDI device listed | Check USB connection; press **f** in the diagnostic to refresh |
| Sidetone, but no events | Select the adapter and press **Enter** to enable MIDI; confirm event count increases |
| “Release both pads” | Release all touches; check held-note count in the diagnostic |
| `[gap?]` or an uncertain symbol | Match sending WPM to adapter speed; inspect gaps with **f**, then replay |
| Decode differs from what you intended | Listen to **y** first; the decoder never fills in the expected answer |
| No computer audio | Check macOS output/volume and restart the app; adapter sidetone is separate |

</details>

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
| o | Operating reference | Operating reference; Space pauses/resumes |
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

## Operating reference: NATO alphabet and RST

Press **o** for operating references. **h/l**, left/right arrows, or **Tab** switches
between the NATO/ITU alphabet and RST scales. **j/k** or up/down arrows scroll; **Esc/o** returns
without changing your practice selection. From browsing, **i** inside the
reference imports just that editable list after confirmation. Top-level **i**
continues to import all starter lists. Existing edits and deletions are preserved.

The phonetic list sends each individual letter in CW and uses its standardized
word as an editable spoken prompt: `A` CW → “Alfa” → `A` CW. It defaults to
`sandwich`; all display, sequence and speech controls still apply. This does not
change the normal letter-name speech for callsigns elsewhere.

The RST list includes `599`, `579`, `339`, `449`, and `567`, with spoken prompts
explaining each Readability–Strength–Tone rating. A KE2ISH/W1AW exchange uses the
existing two-station tones; set line advancement to `automatic` to run it in order.
CW uses three ratings; voice omits tone. These clean synthesized examples teach
report sounds and meanings—they do not simulate the reported signal's strength,
noise or tone quality. Both references follow [ARRL's operating aids](https://www.arrl.org/quick-reference-operating-aids).

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
After capture, tentative decoding and gap ranges appear automatically. **- / +**
sets the sending speed used for interpretation (initially 15 WPM, range 5–60),
independently of listening settings; it does not change adapter speed. **f** opens
full feedback with **j/k** scrolling. Hidden practice waits for **R** before
showing the target or comparison. Unknown marks and ambiguous gaps stay uncertain;
prosign aliases such as `<AR>/+` share one continuous pattern. There is no accuracy
score. The adapter generates element lengths; your spacing is the useful manual
signal. Tolerances are documented in [sending design](docs/design/sending-practice.md).
The live `A` capture/decode/replay walkthrough passed on 2026-10-17.
Unplug/reconnect also passed after reselecting QT Py M0. Broader live feedback
acceptance remains open.

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

## Development and verification

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
implementation, NATO/RST slices 11–12, and Vail slices 13–15. See also the [UI design](docs/design/terminal-ui.md)
and [architecture](docs/design/playback-architecture.md).
