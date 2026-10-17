# Vail input and send-back practice

The user authorized continuing the Vail workflow after the cooperative hardware
probe. This extends the existing session owner and output cleanup barrier.
NATO/RST tickets remain independent queued work.

## Connection

From browsing, m opens a MIDI diagnostic. Select the actual adapter and Enter
to enable USB MIDI. The app pairs the selected source with exactly one output
on the same MIDI entity; it does not guess across devices. Only CC0 is sent:
MIDI while connected, keyboard output on explicit disconnect or orderly exit.
Iambic B, speed, tone and firmware are not changed or queried. Disconnect/reconnect
requires a fresh explicit connection; no adapter is silently reconfigured.
Esc leaves the diagnostic with the connection active for send-back practice.

CoreMIDI callbacks enqueue bounded timestamped packets in native storage. The
Go boundary drains them without callbacks into Go, decodes channel-one notes
0/1/2, and supports running status and zero-velocity releases. Endpoint removal
or queue overflow invalidates the connection. Timestamps use the same monotonic
host clock as capture arming; UI poll time is not the keying clock.

## Exercise

From an ordinary entry or an exchange line, t starts send-back practice.
Reference CW plays once, ignoring automatic sequence/reinforcement policies
for this exercise. At Ready, Enter explicitly arms capture; Enter again finishes
it. Finishing records a host-clock cutoff, allows 100 ms for queued MIDI delivery,
then drains a fresh serialized poll before completing. Events after the cutoff
are excluded. Exceptionally late events from before the cutoff invalidate the
attempt and cancel replay; the settling interval is not a hard latency guarantee.
y replays the attempt, v replays reference, and n starts a fresh reference
and attempt. R reveal and h/l or arrows retain their existing meanings. Settings
and key navigation must retain the session's generation/cancellation protection.

The adapter supplies immediate sidetone; the app produces no duplicate live
sidetone. Attempts retain mark and gap durations; replay synthesizes those exact
durations through existing PCM output rather than quantizing them into text.
Leading/trailing idle silence is omitted, internal gaps retained. Attempts are
session-local, bounded to two minutes and 1024 marks. Ending with a held mark,
unmatched/overlapping notes, reversed timestamps, overflow or disconnect makes
an attempt incomplete and unreplayable. A new capture requires released pads.
Old events preceding capture arming cannot enter the attempt. Pause during
capture abandons it and resumes at Ready; pause during audio restarts that
portion. Navigation clears the attempt. Back/quit stop output before returning.

Ticket 15's decoding/feedback remains a separate dependent slice. No microphone,
radio output, persistent recording, or automatic grading is added here.

## Sending feedback (ticket 15)

After capture, show a tentative decode, timing ranges and a comparison. f opens
full feedback with j/k scrolling; Esc closes it. Long attempts remain inspectable. Sending
WPM is session-local, initially 15, adjusted with -/+ (5–60); it is independent
of listening/Farnsworth settings and does not reconfigure the adapter. Recompute
feedback from the original marks when adjusted. No expected text enters decoding.

Classify durations in units of 1200 / sending WPM milliseconds. Accept dits and
intra-element gaps at 0.65–1.35 units, dahs and character gaps at 2.4–3.6, word
gaps at 5.6–8.4. All other durations are uncertain, not rounded to a nearest
category. These deliberately conservative heuristic windows include the verified
80/240 ms adapter output at 15 WPM; they are not calibrated accuracy guarantees.
Unknown elements use ?, ambiguous boundaries use [gap?]. Continuous patterns may
have multiple letter/punctuation/prosign meanings: display alternatives and
compare code groups and word boundaries, not a guessed semantic label.

Comparison reports matching patterns, definite mismatch, or uncertain timing;
no input and incomplete captures are separate states. Hidden practice may show
one's own sending but suppresses target comparison until R reveal. Timing ranges
show element, character and word gaps in dit units; ambiguous gaps are identified.
The adapter generates element lengths; feedback does not grade raw paddle touch,
radio reception or an accuracy percentage. USB timestamp and human hardware
acceptance limitations from tickets 13/14 remain explicit.
