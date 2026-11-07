# D16: an effect removed while it acts comes back (client evidence for branch `entity-effect-tick-readd`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28. The fix is branch `entity-effect-tick-readd` (commit 3e792558, `server/entity: Do not re-add an effect removed during
its own tick`); this folder is a client reproduction run against that commit unchanged. The code-level walk-through is
in `../entity-effect-tick-readd/`.

## Cause

`EffectManager.Tick` (`server/entity/effect.go`) calls `eff.Type().Apply(entity, eff)` and then writes
`m.effects[i] = eff.TickDuration()` unconditionally. When the wither damage kills the player, `Player.kill` removes all
effects (`RemoveEffect`, which also tells the client), but the write-back puts the wither back in the map. The player
respawns with it and keeps taking damage, while the client was told the effect is gone.

## Steps (driver scenario `d16`, harness scenario `d16`)

Survival, 127.0.0.1:19402. The harness registers `/wither`, which gives its source wither IV for 60 s, and reports
each second the player's health and effects (`status` events). Cinnabar joins, runs `/wither`, dies of it, presses
Respawn on the death screen, and is watched for 15 s (30 frames, one each 0.5 s).

## Results

| | effects at respawn (`respawn` event) | after respawning | deaths |
| --- | --- | --- | --- |
| before (106fdc28) | wither IV, 43.7 s left | health falls from 20 to 0; dies again about 17 s later | 2 |
| after (3e792558) | none | health stays 20, no effect | 1 |

## Files

- `before.gif`, `after.gif`: the 15 s after respawning (30 captures of 320 px, 0.5 s apart). A black caption bar was
  added above the frames; the frames themselves are unedited. Before: the hearts empty and the death screen comes
  back. After: the hearts stay full.
- `*-before.png` (joined), `*-dead.png` (death screen), `*-respawned.png`, `*-respawned-15s-later.png` (before: dead
  again; after: alive, full health).
- `*-server.jsonl`: harness events (`wither-given`, `death`, `respawn`, `status` each second);
  `*-driver.txt`: the driver's transcript.

Versions: Dragonfly 106fdc28 and 3e792558, Cinnabar 7af298d1 (protocol 2225).
