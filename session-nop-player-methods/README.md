# D32 — session-less player session methods panic

Branch: `session-nop-player-methods` (commit 4b2bebbf, based on upstream/master 106fdc28)
Package: server/session

## Verdict
CONFIRMED on unmodified master (runtime reproduction with genuine stack traces).
On the Nop session: HideHudElement (write to nil map), RemoveScoreboard (nil *atomic.Pointer deref),
SendForm / SendDialogue / CloseDialogue (nil handlers map -> failed type assertion) all panic.
The world recovers the panic and re-raises it to whoever awaits the transaction; on a session goroutine
(synchronous Call) nothing recovers it and the process can stop.

## Who can trigger it
Plugin / server code calling these methods on a session-less player (player.Config Session == nil, NPC-like).

## Reproduction
dragonfly-bugs/repro/d32 (standalone): spawns a session-less player and calls each method inside a world
transaction, printing whether it panicked (surfaced as the transaction's error).
- before (wt/upstream-master): all five [PANIC] (see before.txt).
- after  (wt/session-nop-player-methods): all five [ok] (see after.txt).

## Files
server/session/player.go (HideHudElement, SendDialogue, CloseDialogue, SendForm), server/session/text.go (RemoveScoreboard):
early `if s == Nop { return }`, matching existing methods such as SendGameMode.

## Checks
go build, gofmt, go test ./server/session/, golangci-lint ./server/session/ (0 issues).
