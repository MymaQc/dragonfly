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

- `before.mp4` (0.6 MB, 17.2 s, 172 frames at 10 fps, 1280x792), `after.mp4` (0.6 MB, 13.4 s, 134 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of the stone pickaxe mining diamond ore, then the break (an experience
  orb and the experience bar before; nothing after), recorded in a later run of the same scenario than the PNGs.
  Built from frames captured continuously through Cinnabar's own screenshot tool during the run (48 frames at 2.84
  per second before, 49 at 3.72 after); each frame is shown for as long as it lasted, so playback speed is real time.
  A caption bar was added above the frames (the Dragonfly commit, before or after, and the step being shown); the
  frames themselves are unedited.
- `before.gif` (0.6 MB, 17.2 s, 69 frames at 4 fps, 640x396), `after.gif` (0.7 MB, 13.5 s, 54 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `before-ore-broken.png`, `after-ore-broken.png`: just after the ore broke. Before: a yellow experience orb in the
  gap and the experience bar above the hotbar partly filled. After: only breaking particles.
- `*-ore-before.png`: the three ores before breaking.
- `*-server.jsonl`: harness events; `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).
