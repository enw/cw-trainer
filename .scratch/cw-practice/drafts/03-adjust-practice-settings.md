# 03 — Tune and remember each list's practice settings

**What to build:** Adjust character speed, overall speed, both station tones, and repetition pause in the terminal and retain a different setup for each list.

**Blocked by:** 02 — Create and practice saved message lists.

**Status:** draft — design and skeleton must inform feature readiness

- [ ] Provide discoverable controls for character speed, overall speed, learner/other station tones, and the shared inter-portion pause; display their current values.
- [ ] Start with 20/15 WPM, 440 Hz learner / 600 Hz other station, and a two-second pause. Permit equal speeds for standard timing.
- [ ] Validate finite positive speeds/frequencies, overall speed no greater than character speed, and a nonnegative pause. Document supported bounds.
- [ ] Apply speed/tone edits at the next CW repetition, keeping the current rendering consistent; display pending values clearly when applicable.
- [ ] Save settings per list and restore them after switching lists or restarting.
- [ ] Verify setting changes during playback, invalid values, equal-speed timing, and persistence through the session boundary.

## Design reference

Follow the [ticket 00 UI decisions](../../../docs/design/terminal-ui.md) and
[playback architecture](../../../docs/design/playback-architecture.md).
Design is approved. Skeleton software is implemented, with listening acceptance
pending. This ticket has not started.
