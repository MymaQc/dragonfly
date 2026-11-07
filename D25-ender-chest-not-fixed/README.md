# D25 — ender chest slot broadcast checks the block, not the open window

No branch (not fixed — no reproducible defect; see verdict).
Package: server/session

## Verdict
CONFIRMED divergence at code level; NOT A (reproducible) BUG in normal play.

broadcastEnderChestFunc (server/session/player.go) decides whether to send an ender-chest slot update with:
    if _, ok := tx.Block(*s.openedPos.Load()).(block.EnderChest); ok { s.ViewSlotChange(...) }
i.e. it checks whether the block at the opened position is an ender chest, rather than whether the currently
open window is the player's ender chest inventory. It also reads through the captured setup `tx` while
`!s.inTransaction` (outside that transaction).

The behaviour only differs from "check the open window" if the block at openedPos is replaced while the ender
chest UI is open, and the slot callback runs outside a transaction where it cannot reliably observe such a
replacement. openedPos is initialised to &cube.Pos{} for a real session (session.go), so there is no nil
dereference. There is therefore no reproducible incorrect client behaviour in normal gameplay.

## Recommendation
Low severity / latent smell. A cleaner check would compare s.openedWindow.Load() to s.enderChest instead of
reading the block. Left unfixed because there is no demonstrable defect and the change is behavioural; flag for
maintainer discussion rather than shipping a speculative fix.
