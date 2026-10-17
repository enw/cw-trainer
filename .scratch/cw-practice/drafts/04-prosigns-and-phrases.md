# 04 — Practice continuous prosigns and mixed phrases

**What to build:** Choose a prosign by meaning or enter a phrase combining words, callsigns, and explicit prosigns, then hear it correctly.

**Blocked by:** 02 — Create and practice saved message lists.

**Status:** draft — design and skeleton must inform feature readiness

- [ ] Provide a picker of supported common prosigns with meanings, inserting angle-bracket notation such as <AR>.
- [ ] Play a prosign as one continuous signal without inter-letter spacing; ordinary AR remains separately spaced letters.
- [ ] Support mixed items such as CQ CQ DE KE2ISH <AR>, preserving ordinary letter/word spacing around complete prosigns.
- [ ] Validate the entire item before playback; report unknown/malformed prosigns without partially playing or silently modifying the message.
- [ ] Store meaning and default spoken phrase information for supported prosigns, ready for speech playback.
- [ ] Verify mappings and procedural meanings against authoritative references and record the sources.
- [ ] Independently verify continuous prosign timing and mixed-phrase boundaries, including saved-item round trips.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is authored but awaits user review and skeleton validation; this ticket
has not started.
