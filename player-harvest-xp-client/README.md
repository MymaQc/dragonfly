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

- `before-ore-broken.png`, `after-ore-broken.png`: just after the ore broke. Before: a yellow experience orb in the
  gap and the experience bar above the hotbar partly filled. After: only breaking particles.
- `*-ore-before.png`: the three ores before breaking.
- `*-server.jsonl`: harness events; `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).
