# 11 — Learn the NATO phonetic alphabet alongside CW

**What to build:** Practice each letter's CW sound with its NATO/ITU phonetic word as an editable spoken prompt, and consult the complete alphabet from the terminal.

**Blocked by:** 09 — Curated editable starter library (implemented; reuses speech and display modes from 05/06).

**Status:** implemented — source verified, automated checks and review complete

- [x] Add a separate editable phonetic-alphabet starter list covering A–Z. CW content remains the individual letter; the spoken prompt uses its phonetic word.
- [x] Use the standardized spellings, including Alfa and Juliett, and verify every mapping against the source below.
- [x] Expose a concise A–Z reference in the terminal with arrow/vi navigation and a clear return to practice.
- [x] Reuse whole/highlight/hidden display, speak-now, reinforcement, repeat/cycle/shuffle, and per-list settings. Do not change existing generated callsign speech or other lists.
- [x] Provide a deliberate non-destructive way for existing users to add this starter list; first-use installation includes it without reseeding over later edits or deletions.
- [x] Verify all 26 mappings, spoken prompts, save/reload, hidden-answer concealment, and preservation of an existing library.
- [x] Document how this differs from normal letter-name speech and demonstrate letter CW → phonetic word → letter CW.

## Source and scope

Requested 2026-10-16 during the initial implementation.
[ARRL Quick Reference Operating Aids — International Telecommunication Union Phonetic Alphabet](https://www.arrl.org/quick-reference-operating-aids).
The NATO/ITU spelling alphabet is spoken-word familiarization, not a new Morse encoding.
A global phonetic-callsign speech setting is not assumed by this slice.
