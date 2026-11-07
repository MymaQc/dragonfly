# D16: an effect removed while it acts comes back (client evidence for branch `entity-effect-tick-readd`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28. The fix branch `entity-effect-tick-readd` (commit
3e792558, `server/entity: Do not re-add an effect removed during its own tick`) has its own evidence and PR text in
`evidence/entity-effect-tick-readd/`; this folder only adds a client reproduction of it, run against that commit
unchanged. No PR.md here.

## Cause

`EffectManager.Tick` (`server/entity/effect.go`) calls `eff.Type().Apply(entity, eff)` and then writes
`m.effects[i] = eff.TickDuration()` unconditionally. When the wither damage kills the player, `Player.kill` removes all
effects (`RemoveEffect`, which also tells the client), but the write-back puts the wither back in the map. The player
respawns with it and keeps taking damage, while the client was told the effect is gone.

## Steps (driver scenario `d16`, harness scenario `d16`)

Survival, 127.0.0.1:19402. The harness registers `/wither`, which gives its source wither IV for 60 s, and reports
each second the player's health and effects (`status` events). Cinnabar joins, runs `/wither`, dies of it, presses
Respawn on the death screen, and is watched for 15 s (20 s in the run the videos come from).

## Results

| | effects at respawn (`respawn` event) | after respawning | deaths |
| --- | --- | --- | --- |
| before (106fdc28) | wither IV, 43.7 s left | health falls from 20 to 0; dies again about 17 s later | 2 |
| after (3e792558) | none | health stays 20, no effect | 1 |

## Files

- `before.mp4` (5.4 MB, 43.8 s, 2626 frames at 60 fps, 1280x792), `after.mp4` (5.6 MB, 43.8 s, 2626 frames at 60 fps,
  1280x792): H.264, yuv420p, real-time video of /wither, the death, Respawn, then 20 s after respawning (health falls
  and the player dies again before; full health after), recorded in a later run of the same scenario than the PNGs.
  Recorded by Cinnabar's own recorder (record_start with fps 60, fixed_clock false so that the game follows the wall
  clock as the server does, audio off; hidden window), started before the action and stopped once the result is
  visible. record_stop reported: before, 2626 frames recorded, 0 dropped, 0 filled; after, 2626 frames recorded, 0
  dropped, 0 filled. Then, with ffmpeg only: a caption bar added above the picture (the Dragonfly commit, before or
  after, and the step being shown), and a re-encode (H.264, yuv420p, faststart, same frame rate); the recorded
  picture is not otherwise changed.
- `before.gif` (7.2 MB, 43.8 s, 350 frames at 8 fps, 480x298), `after.gif` (5.1 MB, 43.8 s, 350 frames at 8 fps,
  480x298): made from the MP4s with ffmpeg (palettegen/paletteuse), for display in the PR.
- `*-before.png` (joined), `*-dead.png` (death screen), `*-respawned.png`, `*-respawned-15s-later.png` (before: dead
  again; after: alive, full health).
- `*-server.jsonl`: harness events (`wither-given`, `death`, `respawn`, `status` each second);
  `*-driver.txt`: the driver's transcript.

Versions: Dragonfly 106fdc28 and 3e792558, Cinnabar 7af298d1 (protocol 2225).
