# CW listening practice

Status: ready-for-agent

## Problem Statement

The learner wants to become fluent in Morse code through listening and,
eventually, keying. This first application focuses on listening: recognizing
familiar words, prosigns, the callsign KE2ISH, and exchange lines by their
sound rather than mentally assembling individual letters.

The learner needs repeatable practice with familiar material, optional spoken
reinforcement, controllable timing, and a terminal interface that makes it easy
to stay on one item until ready to move on.

## Solution

Build an interactive macOS terminal application for creating and practicing
saved lists of CW items and short exchanges. Start with a useful editable
library, including KE2ISH. Play synthesized CW through the Mac's selected audio
output and use built-in macOS speech for spoken reinforcement.

Keep three choices independent:

| Setting | Choices | Initial behavior |
| --- | --- | --- |
| Playback style | CW only; CW → speech → CW; CW with occasional speech | CW only, proposed initial selection |
| Text display | Whole item visible; current letter highlighted; hidden until revealed | Whole item visible |
| Sequence | Repeat current item; cycle; shuffle | Repeat current item until Next |

For exchanges, practice one line at a time. Distinguish the learner's station
from the other station by tone. Allow automatic advancement as an alternative
to repeating the current line.

## User Stories

1. As a learner, I want to launch an interactive terminal app, so that I can practice without composing shell commands for every item.
2. As a learner, I want an editable starter library, so that I can begin listening immediately.
3. As a learner, I want separate lists for my callsign, common terms, prosigns, and exchanges, so that I can focus each session.
4. As a learner, I want to add and name lists, so that I can organize my own practice material.
5. As a learner, I want to add, edit, remove, and reorder items, so that the material follows my learning needs.
6. As a learner, I want lists saved automatically, so that my changes survive restarting the app.
7. As a learner, I want each list to remember its practice settings, so that returning to it restores my preferred setup.
8. As a learner, I want KE2ISH treated as one practice item, so that I recognize my whole callsign by sound.
9. As a learner, I want common CW terms played as complete items, so that I learn their familiar sound patterns.
10. As a learner, I want a prosign picker with meanings, so that I can select correctly encoded procedural signals.
11. As a learner, I want prosigns sent continuously, so that I do not learn them as separately spaced letters.
12. As a learner, I want phrases containing text, callsigns, and explicit prosigns, so that I can practice realistic material.
13. As a learner, I want to repeat one item until I move on, so that I control when to introduce new material.
14. As a learner, I want to cycle through a list, so that I can practice its items in order.
15. As a learner, I want shuffled practice, so that recognition does not depend on memorized list order.
16. As a learner, I want CW-only playback, so that I can concentrate on the sound.
17. As a learner, I want CW followed by speech followed by CW, so that I can connect a sound with its identity.
18. As a learner, I want occasional speech after repeated CW, so that I receive reinforcement without speech on every repetition.
19. As a learner, I want to configure the speech interval, so that reinforcement matches my needs.
20. As a learner, I want occasional-speech practice to begin with CW, so that I can try recognizing the item first.
21. As a learner, I want to request speech immediately, so that I can identify an unfamiliar item without changing modes.
22. As a learner, I want editable spoken text for each item or exchange line, so that I control what the voice says.
23. As a learner, I want whole-item text visible, so that I can associate its complete sound with its meaning.
24. As a learner, I want optional current-letter highlighting, so that I can inspect where I am in ordinary text.
25. As a learner, I want text hidden until I reveal it, so that I can test recognition without seeing the answer early.
26. As a learner, I want display and playback settings to be independent, so that spoken reinforcement works with either visible or hidden text.
27. As a learner, I want adjustable Farnsworth timing, so that characters retain their sound while spacing gives me time to recognize them.
28. As a learner, I want adjustable tone frequencies, so that listening is comfortable and stations are distinguishable.
29. As a learner, I want adjustable pauses between playback portions, so that I have time to process what I hear.
30. As a learner, I want to change speed and tone while practicing, so that I can adjust difficulty without leaving the session.
31. As a learner, I want pending speed and tone changes applied on the next repetition, so that a playing item stays consistent.
32. As a learner, I want short realistic exchanges practiced line by line, so that I learn how familiar sounds fit into a conversation.
33. As a learner, I want different tones for the two stations, so that I can hear who is speaking.
34. As a learner, I want optional automatic line advancement, so that I can listen through an exchange in order.
35. As a learner, I want pause, navigation, reveal, and speak-now shortcuts, so that practice remains keyboard driven.
36. As a learner, I want both arrow and vi navigation keys, so that the interface fits my habits.
37. As a learner, I want unsupported input and playback failures explained clearly, so that I can correct them without losing my saved work.

## Implementation Decisions

### Scope and boundaries

- Target macOS first; audio uses the selected system output and speech uses
  built-in macOS speech.
- Keep the library, practice-session behavior, Morse timing/audio generation,
  terminal presentation, and operating-system audio/speech integration as
  distinct responsibilities. Choose the smallest architecture that supports
  testing the session's observable behavior.
- Language, terminal library, audio library, storage format, and packaging are
  not yet selected. Do not interpret this spec as a stack decision.
- Treat recognition of complete sounds as the main design goal. Character
  highlighting is optional; do not make letter decoding the default exercise.

### Practice material

- A saved list contains ordinary items and/or exchanges. An item has CW content
  and an editable spoken phrase. An exchange contains ordered lines; each line
  has CW content, an editable spoken phrase, and a station role.
- KE2ISH is the initial learner callsign. Its spoken default is “K E 2 I S H.”
- Permit longer mixed phrases such as “CQ CQ DE KE2ISH <AR>.” Angle brackets
  explicitly distinguish a prosign from ordinary letters.
- A prosign's component elements form one continuous signal. Ordinary words
  and callsigns retain standard internal character gaps even when practiced
  and displayed as a whole item.
- Provide a picker for supported prosigns, showing their meanings and inserting
  the explicit representation. Reject unknown or malformed tokens with useful
  feedback; never silently omit them from transmitted audio.
- Default spoken phrases read ordinary words naturally, CW abbreviations using
  familiar pronunciations (CQ as “see cue”), callsigns character by character,
  and prosigns by meaning. Full lines combine these conventions. The user may
  override the entire spoken phrase.
- Reordering exchange lines is an editing operation. Practice must preserve
  their conversational order even when exchanges are selected by shuffle.

### Timing and audio

| Setting | Initial value |
| --- | --- |
| Character speed | 20 WPM |
| Overall speed | 15 WPM |
| Learner/item tone | 440 Hz |
| Other station tone | 600 Hz |
| Inter-portion pause | 2 seconds |
| Occasional-speech interval | Every 5 CW repetitions |

- Use Farnsworth timing: preserve character-element timing at character speed
  and extend inter-character/inter-word spacing for the requested overall
  speed. The inter-portion practice pause is separate from Morse spacing.
- Use standard PARIS-based WPM calibration. Verify the timing definition and
  supported Morse/prosign mappings against authoritative references during
  implementation; do not derive timings from wall-clock sleeps in UI code.
- Validate positive finite speeds and frequencies and require overall speed
  to be no greater than character speed. Permit equal speeds for standard
  spacing. Exact supported frequency and speed bounds are implementation
  choices that must be documented and validated consistently.
- One configurable pause separates successive CW repetitions and successive
  portions of CW → speech → CW. Apply the same pause around occasional speech.
  Speech must finish before the next portion's pause begins.
- Audio and visible progress must share a timing basis. Tone generation should
  avoid audible clicks at element boundaries without materially changing timing.
- Speed and tone changes during playback take effect on the next complete CW
  repetition; do not alter an item partway through its current CW rendering.
- Pause and navigation must remain responsive while CW, speech, or a gap is
  running. Do not allow competing speech and CW outputs to overlap.

### Playback and display

- Default sequence behavior repeats the current item or exchange line until
  the learner advances. Cycle and shuffle are also supported.
- CW → speech → CW is a complete practice unit, with the configured pause
  between its three portions.
- Occasional-speech mode begins with CW and speaks after the configured count
  of completed CW repetitions. The initial count is five; make it configurable
  for the session and persist the selected value with the list's settings.
- Whole-item display shows the entire ordinary item or current exchange line.
  Letter mode highlights progress through ordinary characters; a prosign is
  highlighted as a single unit, without implying internal letter spacing.
- Hidden mode conceals the current item's answer until R is pressed. Avoid
  leaking it through the current-item title, list preview, or progress display.
  Deliberate spoken reinforcement remains available in hidden mode.
- Exchanges are practiced line by line, repeating each line by default.
  Optional automatic advancement follows the ordered lines using the correct
  station tone. No separately rendered whole-exchange playback is required.

### Terminal interaction and persistence

- Space pauses/resumes; left/right select previous/next item or exchange line;
  R reveals; S speaks now. Support vi navigation alongside arrows.
- Proposed concrete mapping: h/l mirror left/right and j/k mirror down/up in
  list navigation. Navigation shortcuts must not consume text while editing.
- Provide discoverable controls for list/item editing and practice settings,
  and show current playback, display, sequence, and timing selections.
- Automatically save edits and each list's timing, tones, playback style,
  display choice, sequence choice, pause, and speech interval. Use durable
  writes and surface save failures rather than claiming unsaved work is saved.
- Install starter material on first use without overwriting later user edits.
  Keep user data outside the installed application so upgrades preserve it.

### Starter library

- Separate editable lists for KE2ISH, common CW terms, prosigns, and exchanges.
- Include CQ, DE, K, R, TU, 73, QTH, RST, and QRS among the common terms.
- Include a useful set of common prosigns, with verified encodings, meanings,
  and default spoken phrases. Keep ordinary abbreviations distinct from
  continuous prosigns.
- Include short exchanges covering calling CQ, answering, signal reports,
  name/location, repeats, slower sending, and signing off. Use KE2ISH for the
  learner's station and label the other station consistently.
- Review exchange wording and procedural use against authoritative amateur
  radio references before shipping. The approved scope is these categories;
  the exact library entries and example counterpart callsign remain to be chosen.

## Testing Decisions

The user approved testing primarily through a practice-session boundary with a
controllable clock and captured audio, speech, and display output, plus a real
Mac listening and keyboard check. Tests should exercise externally observable
behavior, not private function structure. There is no existing code or test
suite to reuse.

- Drive sessions using user actions and verify emitted playback/display events:
  all three styles, sequence choices, repetition counting, line advancement,
  speaker tones, reveal behavior, and settings applied at repetition boundaries.
- Verify timing with independently calculated known Morse examples and a
  PARIS calibration case, including standard and Farnsworth spacing. Include
  prosigns and mixed phrases so uninterrupted prosigns cannot regress into
  separately spaced letters.
- Check generated audio duration and frequency through a small real rendering
  integration test; event-only assertions cannot prove actual audio correctness.
- Exercise pause, Next, speak-now, and exit during CW, speech, and silence.
  Verify cancelled output does not continue or overlap the next item.
- Verify list edits/settings survive a save-and-reload round trip and starter
  installation preserves user modifications. Include a meaningful save-error case.
- Verify invalid input fails clearly rather than producing partial or silently
  altered CW. Check speech text is passed as data, never interpreted as a shell command.
- On a real Mac, listen to sample items and both station tones, compare visible
  highlighting to audio, exercise keyboard controls and editing, and restart
  the app to confirm persistence. Record this separately from automated results;
  passing simulated tests does not establish audible quality or terminal usability.

## Out of Scope

- Key/paddle input, keying assessment, or transmit practice.
- Radio connections, hardware keying, and over-the-air transmission.
- Linux/Windows support in this first version.
- Audio-file export, cloud sync, accounts, and online speech services.
- Automatic recognition grading, adaptive curricula, and learning analytics.
- A separate continuous whole-exchange playback mode beyond automatic
  line-by-line advancement.

## Further Notes

The user approved the core requirements summarized above and requested this
spec. Ticket creation and implementation remain separate, user-requested stages.
The ready-for-agent status identifies a usable spec; it does not authorize
skipping that workflow.

The following edge behaviors are proposed defaults rather than individually
confirmed requirements. Review them when breaking the spec into tickets:

- Start a new list in CW-only mode with whole-item text visible.
- Hidden text stays revealed through repetitions of the same item and hides
  again when the current item or exchange line changes.
- Reset the occasional-speech count when selecting a different item/line or
  changing playback style. Count only completed CW repetitions.
- With cycle/shuffle plus occasional speech, complete the configured repetitions
  and spoken reinforcement for the current item before advancing; do not count
  unrelated items toward the speech interval.
- Shuffle visits every top-level item/exchange once per round. Preserve lines
  within an exchange in order. Cycle wraps at the end; both stop on pause/exit.
- Navigating during output cancels it and starts the newly selected item.
  Speak-now temporarily interrupts practice, speaks the current phrase, then
  restarts that item's CW from its beginning after the configured pause.
- If the native speech backend cannot resume mid-utterance, resume a paused
  spoken portion from its beginning. Make the behavior consistent and visible.
- Use a single pause setting initially; separate pauses for different portions
  can be considered later.
