# Commands: `<nil>` sent to players (branch `cmd-nil-error`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28 for `<nil>`, fixed by 7d9704b8. Impact is limited: it
needs a command registered without any `Runnable` (`cmd.New(name, description, aliases)` with no runnable, which
`cmd.New` accepts).

No path sending raw Go error text to players was found in Dragonfly's `server/cmd` and `server/session`: the errors
`Command.Execute` sends are translations (`MessageSyntax`, `MessageUsage`, ...), and the Go errors of packet handlers
(forms, inventory requests) are only logged. That part of the suspicion was not reproduced.

## Cause

`Command.Execute` (`server/cmd/command.go`) keeps `leastErroneous` from the runnables it tried and ends with
`output.Error(leastErroneous)`. With no runnable, `leastErroneous` stays nil, and `Output.Error(nil)` falls through to
`errors.New(fmt.Sprint(nil))`, the text `<nil>`, which `Session.SendCommandOutput` sends.

## Steps (driver scenario `cmd`, harness scenario `cmd`)

The harness registers `cmd.New("empty", "A command without overloads", nil)`. Cinnabar sends `/empty` with its chat.
(The command is not in the AvailableCommands packet, since it has no runnable, but the client still sends it.)

## Results

| | chat after `/empty` |
| --- | --- |
| before (106fdc28) | `<nil>` |
| after (7d9704b8) | Unknown command: empty. Please check that the command exists and that you have permission to use it. |

The unknown command message is the one Dragonfly already sends when none of a command's runnables may be run by the
source (`executeRunnable` returns `MessageUnknown`).

## Files

- `before.mp4` (0.1 MB, 4.9 s, 49 frames at 10 fps, 1280x792), `after.mp4` (0.2 MB, 4.9 s, 49 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of /empty sent, and the server's answer in the chat (<nil> before, the
  unknown command message after), recorded in a later run of the same scenario than the PNGs. Built from frames
  captured continuously through Cinnabar's own screenshot tool during the run (18 frames at 3.83 per second before,
  18 at 3.83 after); each frame is shown for as long as it lasted, so playback speed is real time. A caption bar was
  added above the frames (the Dragonfly commit, before or after, and the step being shown); the frames themselves are
  unedited.
- `before.gif` (0.1 MB, 5.0 s, 20 frames at 4 fps, 640x396), `after.gif` (0.1 MB, 5.0 s, 20 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `before-empty-command-output.png`: the last chat line is `<nil>`.
- `after-empty-command-output.png`: the unknown command message, in red.
- `*-pick-outputs.png`, `*-npcs.png`: earlier part of the same run (`@p`, see `cmd-nearest-player`; not fixed on this
  branch).
- `*-server.jsonl`: harness events (`command` with `empty`); `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Repo checks on the branch: all pass, see `checks/cmd-nil-error.log`.
