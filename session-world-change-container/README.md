# D33: a container window survives a change of world (branch `session-world-change-container`)

Verdict: **CONFIRMED** on unmodified upstream master 106fdc28, fixed by ad1d0236. Reproduced with a move between two
worlds of the same dimension (as multi-world servers do with several overworlds): no `ChangeDimension` packet is
sent, so a client keeps the chest window on screen, as Cinnabar did. A move to another dimension (portal) was not
tested; there the client receives `ChangeDimension`, and whether the official client then keeps the window is not
known from here.

## Cause

On a world change the session (`handleWorldSwitch`, `server/session/session.go`) only changes the dimension, the
position and the chunk loader; `containerOpened`, `openedWindow` and `openedPos` stay as they were. Item stack requests
are then handled in transactions of the new world, but `invByID` still returns `openedWindow`, the inventory of the
chest in the old world, which is changed from the new world's goroutine without that world's lock. When the window is
closed later, `closeCurrentContainer` looks the block up at `openedPos` in the new world, so the chest in the old world
is never told the viewer left.

## Steps (driver scenario `d33`, harness scenario `d33`)

Survival, 127.0.0.1:19402. The harness creates a second overworld (`world.Config{...}.New()`, flat, not saved), puts a
chest with 64 diamonds at -3 -60 2 in the main overworld, with an inventory handler (`chestWatcher`) that reports each
change and the world of the transaction of the player making it. It registers `/away <seconds>`, which, after that
delay, removes the player from its world (`tx.RemoveEntity`) and adds it to the second overworld at the same position
(`AddEntityAt`, then `Teleport`), as a plugin would.

Cinnabar runs `/away 4`, opens the chest, waits for the move (`moved-in` event) and 2.5 s more, then, if the window is
still open, shift-clicks the diamonds out of the chest slot and closes it.

## Results

| | after the move | chest change seen by the handler | chest in the main overworld at shutdown | player |
| --- | --- | --- | --- | --- |
| before (106fdc28) | window still open, 64 diamonds shown | `chest-take` of 64 diamonds, chest in `overworld`, player's transaction in `second-overworld` | empty | 64 diamonds |
| after (ad1d0236) | window closed | none | 64 diamonds | nothing |

## Files

- `before.mp4` (0.4 MB, 12.1 s, 121 frames at 10 fps, 1280x792), `after.mp4` (0.4 MB, 11.0 s, 110 frames at 10 fps,
  1280x792): H.264, yuv420p, real-time video of /away 4, the chest opened, the move to the second overworld with the
  window open, the shift-click out of it (before; after, the window closes), then looking where the chest was,
  recorded in a later run of the same scenario than the PNGs. Built from frames captured continuously through
  Cinnabar's own screenshot tool during the run (44 frames at 3.73 per second before, 41 at 3.83 after); each frame
  is shown for as long as it lasted, so playback speed is real time. A caption bar was added above the frames (the
  Dragonfly commit, before or after, and the step being shown); the frames themselves are unedited.
- `before.gif` (0.4 MB, 12.0 s, 48 frames at 4 fps, 640x396), `after.gif` (0.4 MB, 11.0 s, 44 frames at 4 fps,
  640x396): the same, scaled down for display in the PR.
- `*-chest-open-overworld.png`: the chest window in the main overworld, 64 diamonds.
- `before-after-move-to-second-world.png`: on 106fdc28, 2.5 s after the player is in the second world, the chest
  window is still open and still shows the diamonds. `after-after-move-to-second-world.png`: with the fix, it closed.
- `before-taken-in-second-world.png`: on 106fdc28, the diamonds taken from the window, in the player's inventory.
- `*-second-world-no-chest.png`: looking where the chest is in the main world: there is none in the second world.
- `*-server.jsonl`: harness events (`chest-placed`, `moved-out`, `moved-in`, `chest-take`, `chest-before-close`);
  `*-driver.txt`: the driver's transcript.

Screenshots are unedited captures of Cinnabar (hidden window, 1280x720).

Not covered by the client test: a request handled between the world switch and the session's next chunk update
(at most a tick of the session's 50 ms loop) is refused by the new world check in `invByID`; no client timing reaches
it reliably. The removal of the viewer in the old world is not observable through the public API (the viewer set of a
chest is unexported).

Note: the branch `world-loader-foreign-tx-close` (D34) is in the same area of world changes; the two
branches were not tested together.

Versions: Dragonfly 106fdc28, Cinnabar 7af298d1 (protocol 2225). Repo checks on the branch: all pass, see
`checks/session-world-change-container.log`.
