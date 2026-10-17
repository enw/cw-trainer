# cw-looper

A macOS terminal app for learning familiar CW sounds. The current walking
skeleton contains a saved two-station call/answer exchange with **KE2ISH**.
Each line repeats until you move on; you can pause or request its spoken phrase.

## Run

Requires Go 1.26+ and macOS with `/usr/bin/say` and a configured audio output.
The tested development platform is macOS 26.6.2 on Apple silicon. Other macOS
versions/architectures have not been verified.

```sh
go run ./cmd/cw-looper
```

Or build a binary:

```sh
go build -o bin/cw-looper ./cmd/cw-looper
./bin/cw-looper
```

Use a terminal at least 80 columns × 24 rows. First launch saves a small practice
library. Enter opens the list, then exchange, then starts the selected line.
Browsing does not play sound.

| Control | Browse | Practice |
| --- | --- | --- |
| ↑/↓ or k/j | Select | No action |
| Enter or l | Open/start | Enter retries a failed output |
| ←/→ or h/l | Back/open | Previous/next line; stops at the ends |
| [ / ] | — | Previous/next exchange; stops at the ends |
| Space | — | Pause/resume |
| s | — | Speak current phrase |
| Esc | Back | Stop and return to line selection |
| ? | Help | Help (playback continues) |
| q or Ctrl+C | Quit | Stop output and quit |

Pause restarts interrupted CW/speech from the beginning of that portion; a
paused gap resumes its remaining time. Navigation stays silent while paused.
Speak-now while paused returns to paused; while running, it speaks and then
restarts CW after the gap. Switching lines cancels the old output first.

This slice uses fixed **20 WPM character / 15 WPM overall** Farnsworth timing,
**440 Hz for You / 600 Hz for Other**, and **2-second practice gaps**. It shows
the complete current line. Editors, configurable settings, prosigns, alternate
displays, reinforcement patterns, cycle/shuffle, and automatic exchange
advancement belong to later tickets.

## Saved material

The library is stored at:

```text
~/Library/Application Support/cw-looper/library.json
```

Use `--data /path/to/library.json` for a separate practice library. Only one app
may open a given library at once. The initial W1AW/KE2ISH fixture is example
practice text, not a real contact or recording.

Until the editor slice is implemented, quit the app and edit its JSON to change
exchange titles, lines, station roles (`learner` or `other`), and explicit spoken
phrases. Keep IDs unique and settings at the fixed preset. Supported CW input
is A–Z, 0–9, and spaces; lowercase is played as uppercase. Maximum CW length is
512 bytes and spoken text 1024 bytes. Newline/control characters are rejected.

Existing empty libraries remain empty. Invalid or future-version data is reported
without replacement. Correct the file and relaunch; do not delete your library
as an automatic recovery step. Accepted saves use a temporary file, file sync,
rename, and directory sync. A companion `.lock` file may remain after exit;
the OS releases the lock when the process exits, so do not delete it to bypass
another running instance.

Audio initialization/device failures may require checking the system output and
restarting the app. Speech failures allow CW practice to continue after retrying
or restarting CW with Enter. The speech process is invoked directly with stdin
text; no shell evaluates your phrases. Neither backend requires a network service.
Switching audio devices during a session is not yet verified; restart afterward.

## Verify

```sh
go test ./...
go test -race ./...
go vet ./...
# This deliberately plays tones and speech through the Mac's output:
CW_REAL_AUDIO=1 go test -v ./internal/output -run TestMacOutputLifecycle -count=1
```

Tests cover independently calculated Morse timing, actual PCM duration/frequency,
load/save preservation and locking, session interruption rules, and terminal
navigation. Real audio tests are opt-in. Software stop/reap timings do **not**
prove acoustic silence or sound quality; see the
[verification record](docs/design/skeleton-verification.md) for observed results
and outstanding listening checks.

The [dependency map](.scratch/cw-practice/map.md) is the single implementation
plan. The [UI](docs/design/terminal-ui.md) and
[architecture](docs/design/playback-architecture.md) describe the complete design,
including features not present in the walking skeleton.
