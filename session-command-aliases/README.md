# Commands: aliases enum index in AvailableCommands (branch `session-command-aliases`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28, fixed by 088bd791. Impact is limited: it only happens
when a parameter enum has the type Dragonfly gives an aliases enum, `<command>Aliases`.

## Cause

`BuildAvailableCommands` (`server/session/command.go`) sets `aliasesIndex = uint32(len(enumIndices))` and then
`enumIndices[c.Name()+"Aliases"] = aliasesIndex`, while parameter enums are looked up in the same map by type. With a
parameter enum of type `tpAliases`:

- if `/tp` is walked first, the parameter finds `tpAliases` in the map and uses tp's aliases as its options;
- if the parameter enum is walked first, tp's aliases enum replaces its entry without adding a key, so
  `len(enumIndices)` stays one behind `len(enums)` and every later command's aliases offset points at the enum before
  its own.

Commands are walked in map order, which changes from one build to the next.

## Steps (driver scenario `aliases`, harness scenario `cmd`)

The harness registers `first` (alias `f`) with a parameter of an enum type `tpAliases` (options `alpha`, `beta`), `tp`
(alias `teleport`) and `conf` (alias `c`), plus `pick` and `empty` (other bugs). Note that `cmd.New` adds the name of a
command to its aliases.

1. Deterministic check: when a player joins, the harness builds the packet 500 times with
   `session.BuildAvailableCommands(cmd.Commands(), player, nil)` and counts the aliases of each command and the
   options of `/first <which>` (event `alias-check`).
2. On the client: Cinnabar joins four times; each time it opens the chat, types `/` and then `/first `, and the
   suggestions it draws from the packet are captured.

## Results

| | `/conf` aliases wrong (`[teleport tp]`) | `/first <which>` options wrong (`[teleport tp]`) | joins whose `/` suggestions lack `/c` |
| --- | --- | --- | --- |
| before (106fdc28) | 1003 of 2000 builds | 507 of 2000 builds | 2 of 4 (joins 1, 2) |
| after (088bd791) | 0 of 2000 | 0 of 2000 | 0 of 4 |

In an earlier run on 106fdc28 (`before-earlier-run-suggest-first.png`, same commands, run 20261107-030018), the
client showed `/first <teleport|tp>` with the suggestions `teleport` and `tp` instead of `alpha` and `beta`.

## Files

- `before.mp4` (0.4 MB, 34.4 s, 344 frames at 10 fps, 1280x792), `after.mp4` (0.4 MB, 33.9 s, 339 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of four joins (cut between joins), each showing the suggestions for /
  and for /first, recorded in a later run of the same scenario than the PNGs. Built from frames captured continuously
  through Cinnabar's own screenshot tool during the run (128 frames at 3.76 per second before, 128 at 3.81 after);
  each frame is shown for as long as it lasted, so playback speed is real time; the time between segments is cut. A
  caption bar was added above the frames (the Dragonfly commit, before or after, and the step being shown); the
  frames themselves are unedited.
- `before.gif` (0.6 MB, 34.2 s, 137 frames at 4 fps, 640x396), `after.gif` (0.6 MB, 34.0 s, 136 frames at 4 fps,
  640x396): the same, scaled down for display in the PR. In the run the videos come from, join 1 before the fix has
  no `/c` in its suggestions for `/`; joins 2 to 4 happened to get a command order without the problem. After the
  fix, all four joins suggest `/c`.
- `before-join{1..4}-suggest-all.png`: suggestions for `/`; joins 1 and 2 have no `/c` (conf got tp's aliases).
- `after-join{1..4}-suggest-all.png`: `/c` is suggested in every join.
- `*-join{1..4}-suggest-first.png`: suggestions for `/first ` (`alpha`, `beta` in all eight captures).
- `before-earlier-run-suggest-first.png`: `/first <teleport|tp>` on 106fdc28, earlier run.
- `*-server.jsonl`: harness events, with the four `alias-check` counts; `before-earlier-run-server.jsonl`: the events
  of the earlier run.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720). Not tested with the official client: with
the fix, a parameter enum named like an aliases enum keeps its own enum, so the packet then holds two enums with the
same type name (enums are referenced by index). Without such a collision, the packet is unchanged by the fix.

Repo checks on the branch: all pass, see `checks/session-command-aliases.log`.
