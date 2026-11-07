# Ore experience with a tool that cannot harvest it (client evidence for branch `player-harvest-xp`, PR #1486)

Extra evidence for the existing branch `player-harvest-xp` (d5e21503), already proposed upstream as #1486. The branch
was not changed: it was checked out read-only in a detached worktree (`wt/xp-check`) to build the harness.

## Steps (driver scenario `xp`, harness scenario `xp`)

Survival, 127.0.0.1:19402. Three diamond ore blocks at x 2-4, y -60, z 3; the player gets a stone pickaxe (slot 0) and
an iron pickaxe (slot 1). Cinnabar walks to 3.5 -60 1, holds the stone pickaxe on the diamond ore at 3 -60 3 until the
server reports it broken, and looks at where it was.

## Results (harness `block-break` event)

| | drops | xp | client |
| --- | --- | --- | --- |
| before (106fdc28) | none | 4 | an experience orb at the ore's place, the experience bar filling |
| after (d5e21503) | none | 0 | no orb, empty experience bar |

## Files

- `before.mp4` (1.5 MB, 12.6 s, 754 frames at 60 fps, 1280x792), `after.mp4` (1.5 MB, 13.2 s, 791 frames at 60 fps,
  1280x792): H.264, yuv420p, real-time video of the stone pickaxe mining diamond ore, then the break (an experience
  orb and the experience bar before; nothing after), recorded in a later run of the same scenario than the PNGs.
  Recorded by Cinnabar's own recorder (record_start with fps 60, fixed_clock false so that the game follows the wall
  clock as the server does, audio off; hidden window), started before the action and stopped once the result is
  visible. record_stop reported: before, 754 frames recorded, 0 dropped, 0 filled; after, 791 frames recorded, 0
  dropped, 0 filled. Then, with ffmpeg only: a caption bar added above the picture (the Dragonfly commit, before or
  after, and the step being shown), and a re-encode (H.264, yuv420p, faststart, same frame rate); the recorded
  picture is not otherwise changed.
- `before.gif` (3.7 MB, 12.6 s, 251 frames at 20 fps, 640x396), `after.gif` (3.6 MB, 13.2 s, 264 frames at 20 fps,
  640x396): made from the MP4s with ffmpeg (palettegen/paletteuse), for display in the PR.
- `before-ore-broken.png`, `after-ore-broken.png`: just after the ore broke. Before: a yellow experience orb in the
  gap and the experience bar above the hotbar partly filled. After: only breaking particles.
- `*-ore-before.png`: the three ores before breaking.
- `*-server.jsonl`: harness events; `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).
