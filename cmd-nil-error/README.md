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

- `before-empty-command-output.png`: the last chat line is `<nil>`.
- `after-empty-command-output.png`: the unknown command message, in red.
- `*-pick-outputs.png`, `*-npcs.png`: earlier part of the same run (`@p`, see `cmd-nearest-player`; not fixed on this
  branch).
- `*-server.jsonl`: harness events (`command` with `empty`); `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Repo checks on the branch: all pass, see `checks/cmd-nil-error.log`.
