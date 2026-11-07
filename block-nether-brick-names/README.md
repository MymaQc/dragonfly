# D13: nether brick slab and stairs names swapped (branch `block-nether-brick-names`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28, fixed by b0233a05.

## Cause

`encodeSlabBlock` (`server/block/slab_type.go`) and `encodeStairsBlock` (`server/block/stairs_type.go`) return
`nether_brick` when `block.Type == RedNetherBricks()` and `red_nether_brick` otherwise. Walls (`wall_type.go`) return
the right names. The swap dates back to the first implementation of these slabs (4080c9b3, 2022: red nether bricks
encoded as `stone_slab` type 7, the game's nether brick slab).

The game's names: in the resource pack of Mojang's Bedrock Dedicated Server 1.26.62.1
(`resource_packs/vanilla_1.21.60/blocks.json`), `red_nether_brick_slab` uses the texture `red_nether_brick_slab` and
`nether_brick_slab` the texture `nether_brick_slab`; `vanilla/texts/en_US.lang` names
`tile.red_nether_brick_stairs` "Red Nether Brick Stairs" and `tile.nether_brick_stairs` "Nether Brick Stairs". The
client below shows the same: its own name and texture for each identifier.

Since encoding and decoding are swapped alike, saved worlds and the client always use the game's names; the Go values
are the wrong way round, for any plugin or code that builds these slabs and stairs.

## Steps (driver scenario `d13`, harness scenario `d13`)

Creative, 127.0.0.1:19402. The harness places, at y -60, z 3, from x 1 to 5: `Slab{Block: NetherBricks{}}` at x 1,
`Stairs{Block: NetherBricks{}}` at x 2, `Slab{Block: NetherBricks{Type: RedNetherBricks()}}` at x 4 and
`Stairs{Block: NetherBricks{Type: RedNetherBricks()}}` at x 5; behind each, at z 4, the full block asked for
(`NetherBricks{}` behind x 1-2, `NetherBricks{Type: RedNetherBricks()}` behind x 4-5). The joining player gets, in
hotbar slots 0-5: normal slab, red slab, normal stairs, red stairs, normal block, red block (16 each).

Cinnabar walks to 3.5 -60 0 and looks at the row (towards +Z, so **x grows to the left of the image**: the red asks
are on the left, the normal ones on the right), then selects hotbar slot 1 (the slab asked as red nether bricks).

## Results

| hotbar slot, asked as | identifier the client reports, before | after |
| --- | --- | --- |
| 0, normal nether brick slab | minecraft:red_nether_brick_slab | minecraft:nether_brick_slab |
| 1, red nether brick slab | minecraft:nether_brick_slab | minecraft:red_nether_brick_slab |
| 2, normal nether brick stairs | minecraft:red_nether_brick_stairs | minecraft:nether_brick_stairs |
| 3, red nether brick stairs | minecraft:nether_brick_stairs | minecraft:red_nether_brick_stairs |
| 4, normal nether bricks | minecraft:nether_brick | minecraft:nether_brick |
| 5, red nether bricks | minecraft:red_nether_brick | minecraft:red_nether_brick |

## Files

- `before.mp4` (0.9 MB, 12.7 s, 762 frames at 60 fps, 1280x792), `after.mp4` (0.9 MB, 12.7 s, 762 frames at 60 fps,
  1280x792): H.264, yuv420p, real-time video of the row of slabs and stairs, then hotbar slot 2 (slab asked as red)
  and slot 1 (slab asked as normal) held, with the client's item name, recorded in a later run of the same scenario
  than the PNGs. Recorded by Cinnabar's own recorder (record_start with fps 60, fixed_clock false so that the game
  follows the wall clock as the server does, audio off; hidden window), started before the action and stopped once
  the result is visible. record_stop reported: before, 762 frames recorded, 0 dropped, 0 filled; after, 762 frames
  recorded, 0 dropped, 0 filled. Then, with ffmpeg only: a caption bar added above the picture (the Dragonfly commit,
  before or after, and the step being shown), and a re-encode (H.264, yuv420p, faststart, same frame rate); the
  recorded picture is not otherwise changed.
- `before.gif` (2.5 MB, 12.7 s, 254 frames at 20 fps, 640x396), `after.gif` (2.4 MB, 12.7 s, 254 frames at 20 fps,
  640x396): made from the MP4s with ffmpeg (palettegen/paletteuse), for display in the PR.
- `before-row.png`: on 106fdc28, the slab and stairs asked as red (left) are dark in front of a red block; those asked
  as normal (right) are red in front of a dark block.
- `after-row.png`: with the fix, each slab and stairs has the colour of the full block behind it.
- `before-held-slot1-asked-red-slab.png`: slot 1 (asked as red) selected; the client names it "Nether Brick Slab" and
  draws it dark. `after-held-slot1-asked-red-slab.png`: "Red Nether Brick Slab", red.
- `*-inventory.png`: the creative inventory with the hotbar (not needed for the verdict).
- `*-server.jsonl`: harness events (`placed` and `kit`, with the Go value asked and the name it was encoded as);
  `*-driver.txt`: the driver's transcript, with the identifiers the client reports.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Versions: Dragonfly 106fdc28 (protocol 2225 / 1.26.60), Cinnabar 7af298d1 (protocol 2225), BDS 1.26.62.1 data.
Repo checks on the branch: all pass, see `checks/block-nether-brick-names.log`.
