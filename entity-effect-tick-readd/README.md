# D16 — an effect removed during its own tick is re-added

Branch: `entity-effect-tick-readd` (commit 3e792558, based on upstream/master 106fdc28)
Package: server/entity

## Verdict
CONFIRMED at code level. The live reproduction with a client is in `../entity-effect-tick-readd-client/`.
EffectManager.Tick applies each lasting effect then writes `m.effects[i] = eff.TickDuration()` unconditionally.
If Apply removes the effect (or replaces it), the write-back re-adds the removed effect or overwrites the
replacement. The write-back now runs only when the same effect is still present after Apply.

## Who can trigger it
An effect whose Apply removes effects during its own tick (gameplay interaction or plugin code).

## Files
server/entity/effect.go (Tick: guard the duration write-back with `if cur, ok := m.effects[i]; ok && cur == eff`).

## Checks
go build, gofmt, go test ./server/entity/, golangci-lint ./server/entity/ (0 issues).
No standalone program here: mocking the large entity.Living interface is not warranted, and the bug is reproduced live in `../entity-effect-tick-readd-client/`.
