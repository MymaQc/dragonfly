# D17 — unknown effect id in saved data keeps a zero-value effect

Branch: `player-effect-unknown-id` (commit 7d598ed9, based on upstream/master 106fdc28)
Package: server/player/playerdb

## Verdict
CONFIRMED (reproduced with a playerdb test).
dataToEffects allocated the result slice at len(data) and indexed into it, `continue`-ing on an unknown effect
id. The skipped index kept a zero-value effect.Effect (nil Type), which the player then fails on when adding it.

## Who can trigger it
World / player data only: saved player data containing an effect with an id the server does not know
(version mismatch, corruption, hand-edited data).

## Reproduction
Test kept OUTSIDE the fix commit: evidence/player-effect-unknown-id/d17_repro_test.go.txt
Copy to server/player/playerdb/, `go test ./server/player/playerdb/ -run TestD17 -v`.
- before (wt/upstream-master): FAIL — "effect 0 has a nil Type (zero-value effect kept for unknown id)" (before.txt)
- after  (wt/player-effect-unknown-id): PASS — 0 effects returned (after.txt)

## Files
server/player/playerdb/effect.go (build the slice with append, skipping unknown ids).

## Checks
go build, gofmt, go test ./server/player/playerdb/, golangci-lint ./server/player/playerdb/ (0 issues).
