# Vail touch-pad MIDI verification — 2026-10-16

## Result

The user's actual adapter successfully sent USB MIDI keying events to this Mac.
The user used its capacitive pads and reported hearing local sidetone. Their
reported configured keyer is Iambic B; the probe did not change or independently
query that setting. No microphone or audio-input capture was used.

The standalone diagnostic used macOS CoreMIDI, outside the app. The sandbox
could not create a MIDI client (status -10833); the authorized host probe
succeeded (status 0). There was exactly one source and destination named
`QT Py M0`. Exact installed firmware/version remains unknown.

## Procedure and observations

1. User connected the adapter fresh. The probe selected its named endpoint pair.
2. Connected an input listener, then sent `B0 00 00` (CC0 MIDI mode). The send
   returned success and an incoming matching CC0 packet was observed.
3. User tapped dit three times, then dah three times:

| Mark | Down timestamp (s) | Up timestamp (s) | Duration (ms) |
| --- | ---: | ---: | ---: |
| Dit | 6.949423 | 7.029218 | 79.795 |
| Dit | 8.291394 | 8.371500 | 80.106 |
| Dit | 9.761593 | 9.841288 | 79.695 |
| Dah | 11.319653 | 11.559292 | 239.639 |
| Dah | 12.725676 | 12.965304 | 239.628 |
| Dah | 14.438734 | 14.678628 | 239.894 |

4. Holding the dit pad yielded repeated ~80 ms marks separated by ~80 ms gaps.
   Holding the dah pad yielded repeated ~240 ms marks with ~80 ms gaps.
5. Squeezing both pads yielded alternating dit/dah marks. All observed downs
   had matching ups; the final release was at 52.040550 s and no further keying
   events arrived before the diagnostic was ended around 61.8 s.
6. Sent `B0 00 7F` to return to keyboard output. Send status was success and the
   incoming matching CC0 packet was observed. Ports/client closed and the probe
   exited successfully. No speed, tone, keyer-program or firmware writes occurred.

Times are CoreMIDI host packet timestamps relative to probe startup, not acoustic
measurements. ~80 ms dits imply ~15 WPM using 1200/dit_ms. This proves event
capture and nominal timing, not end-to-end latency or touch-to-sound latency.

## Important protocol difference

Observed dit packets: `90 01 7F` / `80 01 00`.
Observed dah packets: `90 02 7F` / `80 02 00`.

The [current upstream specification](https://github.com/Vail-CW/vail-adapter/blob/master/docs/MIDI_INTEGRATION_SPEC.md)
describes note 0 for onboard-keyer marks and notes 1/2 for passthrough/direct
routing. This adapter instead produced regularly timed repetitions using notes
1/2 while the user reported Iambic B. Do not infer passthrough solely from note
numbers or assume note 0 is the only valid timed output. Integration must use
observed down/up timing and explicitly support the tested device behavior.
The squeeze trace alone does not independently distinguish Iambic A from B.

## Integrated implementation and remaining acceptance

Tickets 13–15 now provide a native MIDI diagnostic, reference/capture/replay,
and tentative decoding with spacing feedback. Automated checks cover recorded
mark lengths, event handling, capture interruption, replay timing, independent
sending speed, uncertainty, prosigns and hidden-target comparison. The app has
enumerated and connected to QT Py M0 and played a reference. On 2026-10-17, the user reported running the README A capture/decode/replay
walkthrough with no problems. This accepts the basic live send-back path, including
replay listening and A feedback. No new timing measurements were collected.
The user also confirmed recovery after unplug/reconnect by explicitly reselecting
QT Py M0 on 2026-10-17. Broader live feedback cases (word spacing, prosigns,
uncertain or mismatched attempts) remain unconfirmed. The temporary standalone Swift probes were used
from `/tmp`; they are not installed features or repository source.
