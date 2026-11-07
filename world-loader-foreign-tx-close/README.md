# D34 — Loader closed with another world's transaction

Branch: `world-loader-foreign-tx-close` (commit d9d49876, based on upstream/master 106fdc28)
Package: server/world

## Verdict
CONFIRMED (runtime reproduction).
Loader.Close removed the loader's chunk viewers through tx.World() (the transaction's world). removeViewer keys
off the chunks of the world it is called on, so when Close is given a transaction of another world (a session
closed while changing worlds), it looked the chunks up in the wrong world and left the loader viewing every
chunk of its own world, never hiding their entities. evictUnused and the rest of the file already use the
loader's own world (l.w); Close was the outlier.

## Who can trigger it
Internal: a session closed during a world change (player left while travelling). Server/plugin code closing a
loader with a foreign transaction.

## Reproduction
Test kept OUTSIDE the fix commit: evidence/world-loader-foreign-tx-close/d34_repro_test.go.txt
Loads a loader in world A, closes it with world B's transaction, counts A's chunks still viewed by the loader.
- before (wt/upstream-master): FAIL — "loader still views 69 chunks of world A after close with a foreign tx" (before.txt)
- after  (wt/world-loader-foreign-tx-close): PASS — 0 chunks (after.txt)

## Files
server/world/loader.go (Close removes viewers through l.w; when the transaction is for another world it schedules
the removal on l.w via l.w.exec, mirroring ChangeWorld, so a transaction never touches another world).

## Checks
go build, gofmt, go test ./server/world/ (full suite passes), golangci-lint ./server/world/ (0 issues).
