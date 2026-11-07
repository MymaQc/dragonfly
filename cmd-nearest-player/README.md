# Commands: `@p` does not select the nearest player (branch `cmd-nearest-player`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28, fixed by 8ab74761.

## Cause

`parser.parseTargets` (`server/cmd/argument.go`), case `@p`, fills `playerDistances[i]` for each player, then calls
`sort.Slice(players, func(i, j int) bool { return playerDistances[i] < playerDistances[j] })`. `sort.Slice` swaps the
elements of `players` but not of `playerDistances`, so after the first swap the comparison reads the distances of
other players. The result depends on the order of `tx.Players()` (map order) and is often wrong with three or more
players.

## Steps (driver scenario `cmd`, harness scenario `cmd`)

Creative, 127.0.0.1:19402. The harness adds three players without a session (as NPCs), with plain skins, 2 blocks
south of the spawn: `Nearest` (blue) 1.5 blocks east, `Near` (green) 3 blocks east, `Far` (red) 6 blocks east. It
registers `/pick <targets>`, which prints `selected <name> (<distance> blocks)` and reports every player's distance
(`pick` event, `players_nearest_first`). Cinnabar (the source, distance 0) runs `/pick @p` 20 times.

In the game, `@p` run by a player selects that player, the nearest at distance 0. Whatever one thinks of that, a
selection of `Near` (3.6) while `Nearest` (2.5) is closer, as seen in earlier runs, can be right under no reading.

## Results (20 runs each)

| | Cinnabar (0.0) | Nearest (2.5) | Near (3.6) |
| --- | --- | --- | --- |
| before (106fdc28) | 9 | 11 | 0 |
| after (8ab74761) | 20 | 0 | 0 |

Earlier runs on 106fdc28 with the same setup: 3 / 15 / 2 (run 20261107-030203-cmd-base), and with only Near and Far:
3 Cinnabar, 17 Near (run 20261107-030018).

## Files

- `before.mp4` (0.2 MB, 18.2 s, 182 frames at 10 fps, 1280x792), `after.mp4` (0.2 MB, 18.5 s, 185 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of the three NPCs, then /pick @p run 20 times, with each selected player
  printed in the chat, recorded in a later run of the same scenario than the PNGs. Built from frames captured
  continuously through Cinnabar's own screenshot tool during the run (68 frames at 3.81 per second before, 69 at 3.78
  after); each frame is shown for as long as it lasted, so playback speed is real time. A caption bar was added above
  the frames (the Dragonfly commit, before or after, and the step being shown); the frames themselves are unedited.
- `before.gif` (0.3 MB, 18.2 s, 73 frames at 4 fps, 640x396), `after.gif` (0.3 MB, 18.5 s, 74 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `before-pick-outputs.png`: the chat after the 20 runs on 106fdc28 (last 10 lines): `Nearest` and `Cinnabar`
  alternate. `after-pick-outputs.png`: `Cinnabar` every time.
- `*-npcs.png`: the three NPCs as the client sees them.
- `*-empty-command-output.png`: the same run continued with `/empty` (see `cmd-nil-error`; on this branch `<nil>` is
  still sent).
- `*-server.jsonl`: harness events (`command`, and `pick` with the targets and all distances);
  `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Versions: Dragonfly 106fdc28, Cinnabar 7af298d1 (protocol 2225). Repo checks on the branch: all pass, see
`checks/cmd-nearest-player.log`.
