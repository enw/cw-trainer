# 15 — Compare sent text and spacing with the practice target

**What to build:** After a captured send-back attempt, show a tentative Morse decode and useful timing/spacing comparisons alongside replay, helping the learner improve without overstating accuracy.

**Blocked by:** 14 — Listen/send/replay with verified captured timing.

**Status:** draft — queued follow-up; implementation not requested

- [ ] Decode captured marks/gaps with an explicit sending-speed estimate or setting, independently of the reference's listening/Farnsworth speeds.
- [ ] Display uncertain elements and boundaries as uncertain; distinguish no input, incomplete capture, ambiguous timing and a confident mismatch.
- [ ] Compare the attempt with the current item/line, including continuous prosigns and word boundaries, without silently filling missing input from the expected answer.
- [ ] Offer practical timing information, especially character and word spacing, with replay of the attempt and reference. Explain that onboard Iambic B generates element lengths; those lengths do not measure the learner's raw paddle-contact timing.
- [ ] Keep feedback compatible with hidden practice: do not reveal the target before the learner explicitly requests reveal/comparison.
- [ ] Establish documented tolerances and timestamp limitations from the verified hardware capture. Avoid an unexplained accuracy score or claims of grading radio reception.
- [ ] Test known correct attempts, wrong/extra/missing symbols, uncertain gaps, differing speeds, prosigns, and incomplete/disconnected captures; check representative feedback against real touch-pad attempts.

## Boundary

Exact feedback presentation and tolerance choices need design review before
implementation. Adaptive curricula, persistent performance analytics and
simulated signal-quality grading are outside this slice.
