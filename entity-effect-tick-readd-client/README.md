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

- `before.mp4` (2.3 MB, 44.8 s, 448 frames at 10 fps, 1280x792), `after.mp4` (1.4 MB, 43.5 s, 435 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of /wither, the death, Respawn, then 20 s after respawning (health falls
  and the player dies again before; full health after), recorded in a later run of the same scenario than the PNGs.
  Built from frames captured continuously through Cinnabar's own screenshot tool during the run (166 frames at 3.73
  per second before, 161 at 3.72 after); each frame is shown for as long as it lasted, so playback speed is real
  time. A caption bar was added above the frames (the Dragonfly commit, before or after, and the step being shown);
  the frames themselves are unedited.
- `before.gif` (7.0 MB, 44.8 s, 179 frames at 4 fps, 560x346), `after.gif` (4.9 MB, 43.5 s, 174 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `*-before.png` (joined), `*-dead.png` (death screen), `*-respawned.png`, `*-respawned-15s-later.png` (before: dead
  again; after: alive, full health).
- `*-server.jsonl`: harness events (`wither-given`, `death`, `respawn`, `status` each second);
  `*-driver.txt`: the driver's transcript.

Versions: Dragonfly 106fdc28 and 3e792558, Cinnabar 7af298d1 (protocol 2225).
