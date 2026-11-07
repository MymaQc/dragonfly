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

- `before.mp4` (0.7 MB, 31.7 s, 1900 frames at 60 fps, 1280x792), `after.mp4` (0.7 MB, 31.4 s, 1881 frames at 60 fps,
  1280x792): H.264, yuv420p, real-time video of four joins, each showing the suggestions for / and for /first,
  recorded in a later run of the same scenario than the PNGs. Recorded by Cinnabar's own recorder (record_start with
  fps 60, fixed_clock false so that the game follows the wall clock as the server does, audio off; hidden window),
  started before the action and stopped once the result is visible (four recordings, one per join; the time between
  joins is not recorded). record_stop reported: before, 1900 frames recorded, 0 dropped, 0 filled; after, 1881 frames
  recorded, 0 dropped, 0 filled. Then, with ffmpeg only: a caption bar added above the picture (the Dragonfly commit,
  before or after, and the step being shown), the recordings joined in order, and a re-encode (H.264, yuv420p,
  faststart, same frame rate); the recorded picture is not otherwise changed.
- `before.gif` (2.0 MB, 31.6 s, 633 frames at 20 fps, 640x396), `after.gif` (2.0 MB, 31.4 s, 627 frames at 20 fps,
  640x396): made from the MP4s with ffmpeg (palettegen/paletteuse), for display in the PR.
  In the run the videos come from, joins 2, 3 and 4 before the fix have no `/c` in their suggestions for `/` (join 1
  happened to get a command order without the problem); after the fix, all four joins suggest `/c`.
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
