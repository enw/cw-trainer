# 01 — Hear your callsign in a terminal practice loop

**What to build:** Launch the macOS terminal app and hear KE2ISH repeatedly while seeing the whole callsign, with working pause and exit controls.

**Blocked by:** None — can start immediately after implementation is requested.

**Status:** draft — breakdown awaiting review

- [ ] Choose and document a minimal implementation stack, actual run/test commands, and local installation prerequisites.
- [ ] Play KE2ISH at 440 Hz with 20 WPM character speed, 15 WPM overall Farnsworth timing, and a two-second repetition pause.
- [ ] Display the callsign as a whole item and show playback state and keyboard help.
- [ ] Space pauses/resumes during CW and gaps; exit stops audio and restores the terminal. Document whether an interrupted CW item resumes or restarts.
- [ ] Use standard Morse and PARIS-based timing verified against authoritative references; preserve ordinary internal letter spacing.
- [ ] Keep input responsive during audio output and surface output-device failures clearly.
- [ ] Verify session behavior with a controllable clock and captured outputs, and independently check real rendered tone frequency and duration, including Farnsworth calibration.
- [ ] Perform and record a real Mac playback smoke check; do not claim audible quality solely from automated checks.

## Boundary

This slice uses one fixed practice item. Editing, persistence, speech, alternate display modes, and runtime tuning follow in later slices.
