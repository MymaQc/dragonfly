# D23: container contents lost on save (branch `session-container-save`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28, fixed by f1928f11.

## Cause

`World.saveChunk` only stores a chunk whose `modified` flag is set. The flag is set by block, biome and entity
changes, but changing the inventory of a container (`Inventory.SetItem`, called by
`ItemStackRequestHandler.setItemInSlot` for a client's request on the open window) does not set it. The flag is never
cleared after a save, so the loss happens when nothing else changed in that chunk since it was loaded. A player
joining marks the chunk it spawns in, so the chest is placed in another chunk.

Not covered by the fix: a hopper inserting into a container in another chunk (`Hopper.insertItem` marks only the
hopper's chunk) and the smelting tick of a furnace that stays lit (its inventory changes without `SetBlock`). Neither
was reproduced here.

## Steps (driver scenario `d23`, harness scenarios `d23-place`, `d23`, `d23-check`)

Same world folder for the three servers, survival, 127.0.0.1:19402, flat world, spawn 0 -60 0 (chunk 0, 0).

1. `server1-place`: the harness places a chest at -3 -60 2 (chunk -1, 0) in a world transaction, and the server stops,
   saving it. No client.
2. `server2-fill`: the harness reports the chest loaded from disk (`chest-loaded`, empty) and gives the joining player
   64 diamonds. Cinnabar joins at 0.5 -60 0.5 (chunk 0, 0) without walking, uses the chest, shift-clicks the
   diamonds into it, closes it and quits. Nothing else is done in chunk (-1, 0). Before the server stops, the harness
   reports the chest's contents in memory (`chest-before-close`: 64 diamonds), then `Server.Close` saves the world.
3. `server3-reopen`: the harness reports the chest loaded from disk (`chest-loaded`) before anyone joins; Cinnabar
   joins and opens the chest. No items are given in this step.

## Results

| | server 2, in memory before saving | server 3, loaded from disk | server 3, seen by the client |
| --- | --- | --- | --- |
| before (106fdc28) | 64 diamonds | empty | 0 diamonds |
| after (f1928f11) | 64 diamonds | 64 diamonds | 64 diamonds |

## Files

- `before-server2-chest-filled.png`, `after-server2-chest-filled.png`: step 2, the chest window after the diamonds
  were put in (64 diamonds in the first chest slot).
- `before-server3-chest-reopened.png`: step 3 on 106fdc28, the same chest after the restart, empty.
- `after-server3-chest-reopened.png`: step 3 with the fix, the chest still holds the 64 diamonds.
- `*-server{1,2,3}-*.jsonl`: the harness events (`chest-placed`, `chest-loaded`, `chest-before-close`, joins);
  `*.err.log`: Dragonfly's log; `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Versions: Dragonfly 106fdc28 (gophertunnel v1.63.0, protocol 2225 / 1.26.60), Cinnabar 7af298d1 (protocol 2225),
Go 1.26.5. Repo checks on the branch (gofmt, go test ./..., golangci-lint v2.11.4, go generate ./server/block/ with no
diff): all pass, see `checks/session-container-save.log`.
