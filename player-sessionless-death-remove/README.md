# D15 — a session-less player that dies stays in its world

Branch: `player-sessionless-death-remove` (commit 679e654d, based on upstream/master 106fdc28)
Package: server/player

## Verdict
CONFIRMED on unmodified master (runtime reproduction).
finishDying closes a session-less player; close() routes to respawn when `p.session() != nil`, but session()
never returns nil (it returns session.Nop), so the respawn path runs and immediately bails for a Nop session.
The player is therefore neither respawned nor removed and stays dead in the world forever.

## Who can trigger it
Plugin / server code that creates a session-less player (player.Config with Session == nil, e.g. an NPC) which dies.

## Reproduction
dragonfly-bugs/repro/d15 (standalone): spawns a session-less player, kills it with void damage, waits for
finishDying, then checks whether a player entity is still in the world.
- before (wt/upstream-master): "session-less player is STILL in the world after death (BUG)"
- after  (wt/player-sessionless-death-remove): "session-less player was removed from the world after death (fixed)"

## Files
server/player/player.go (close(): check p.s instead of p.session()).

## Checks
go build, gofmt, go test ./server/player/, golangci-lint ./server/player/ (0 issues).
