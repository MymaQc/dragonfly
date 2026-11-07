# D14: plant drops keep the state of the broken block (branch `block-crop-drops`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28 for the drops Dragonfly hands to `HandleBlockBreak`
(and anything else that uses `BreakInfo().Drops` directly); **not visible** when the drops are left to fall as item
entities, as in plain survival play. Fixed by 954cbf01.

## Cause

`BreakInfo` of these blocks returns the block value itself as the dropped item: `cropSeedDrops(s, ...)` for wheat and
beetroot, `item.NewStack(c, n)` for carrots, potatoes, nether wart and cocoa beans, `oneOf(x)` for pumpkin and melon
stems, sugar cane, kelp, cactus, vines and double flowers, `grassDrops(d)` for double tall grass. The value keeps its
`Growth`, `Age`, attachments or `UpperPart`. The network name is the plain item, so the client cannot tell, but the
`UseOnBlock` of most of these places the item value as it is (`place(tx, pos, c, ...)`), so it plants the broken state.

Why it does not show with item entities: `ItemBehaviourConfig.New` (`server/entity/item_behaviour.go`) does
`item.ReadNBT(item.WriteNBT(i, true), nil)`, which rebuilds the stack from its name and drops the state.

## Steps (driver scenarios `d14-autopickup` and `d14`, harness scenarios of the same names)

Survival, 127.0.0.1:19402. At z 3: water at x 1, hydrated farmland from x 2 to 5 with wheat at growth 7 on x 2 and
carrots at growth 7 on x 4, soul sand at x 6 and 7 with nether wart at age 3 on x 6.

- `d14-autopickup`: the harness's `HandleBlockBreak` adds the drops it is passed to the player's inventory and clears
  them, as auto-pickup plugins do (`auto-pickup` events).
- `d14`: the drops fall as item entities and Cinnabar walks over them to pick them up.

Cinnabar breaks the wheat, the carrots and the nether wart by hand, then plants wheat seeds on the farmland at x 3,
a carrot on x 5 and nether wart on the soul sand at x 7, and looks at the field (towards +Z: **x grows to the left**).

## Results (state the server placed, from its `block-place` events)

| | wheat seeds on x 3 | carrot on x 5 | nether wart on x 7 |
| --- | --- | --- | --- |
| auto-pickup, before (106fdc28) | growth 7 | growth 7 | age 3 |
| auto-pickup, after (954cbf01) | growth 0 | growth 0 | age 0 |
| item entities, before (106fdc28) | growth 0 | growth 0 | age 0 |

## Files

- `before.mp4` (3.1 MB, 43.5 s, 435 frames at 10 fps, 1280x792), `after.mp4` (3.0 MB, 43.2 s, 432 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of the field, each crop broken by hand (drops going to the inventory),
  then the seeds, carrot and nether wart planted (grown at once before, seedlings after), recorded in a later run of
  the same scenario than the PNGs. Built from frames captured continuously through Cinnabar's own screenshot tool
  during the run (154 frames at 3.56 per second before, 154 at 3.59 after); each frame is shown for as long as it
  lasted, so playback speed is real time. A caption bar was added above the frames (the Dragonfly commit, before or
  after, and the step being shown); the frames themselves are unedited.
- `before.gif` (4.3 MB, 43.5 s, 174 frames at 4 fps, 640x396), `after.gif` (4.1 MB, 43.2 s, 173 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `before-field-replanted.png`: auto-pickup on 106fdc28; the wheat, carrots and nether wart just planted are fully
  grown.
- `after-field-replanted.png`: auto-pickup with the fix; they are seedlings.
- `before-field-before.png`, `after-field-before.png`: the field before anything is broken.
- `vanilla-before-*`: the same on 106fdc28 with item entities: seedlings, as the drops were rebuilt from NBT.
- `*-server.jsonl`: harness events (`block-break` with the drops and their Go value, `auto-pickup` / `item-pickup`,
  `block-place` with the properties placed); `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Versions: Dragonfly 106fdc28, Cinnabar 7af298d1 (protocol 2225). Repo checks on the branch: all pass, see
`checks/block-crop-drops.log`.
