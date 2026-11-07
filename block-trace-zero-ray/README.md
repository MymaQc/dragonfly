# D18 — zero-length ray panics in trace.TraverseBlocks

Branch: `block-trace-zero-ray` (commit e5ab2a3a, based on upstream/master 106fdc28)
Package: server/block/cube/trace

## Verdict
CONFIRMED for the direct API: TraverseBlocks(v, v, f) panics ("start and end points are the same").
Explosion reachability: NOT REPRODUCED. The two in-tree projectile callers (entity/projectile.go,
entity/firework_behaviour.go) already guard with a zero-length check before calling trace.Perform. The only
unguarded caller is ExplosionConfig.exposure (block/explosion.go:249), which traces origin -> a point on each
nearby entity's bounding box. Entities exactly at the origin (dist == 0) are skipped, so a panic needs a bbox
sample point to exactly equal the explosion origin (float-exact); I could not produce that in practice. So in
normal gameplay this is essentially unreachable; it is primarily an API-robustness fix.

## Who can trigger it
Direct trace.TraverseBlocks API (plugin code). The explosion path is theoretically reachable but was not
reproduced.

## Reproduction
dragonfly-bugs/repro/d18 (standalone): TraverseBlocks with start == end.
- before (wt/upstream-master): [PANIC] (before.txt)
- after  (wt/block-trace-zero-ray): visits 1 block, no panic (after.txt)

## Files
server/block/cube/trace/trace.go (zero-length ray now visits the start block once).

## Checks
go build, gofmt, go test ./server/block/cube/trace/ and ./server/block/, golangci-lint ./server/block/cube/trace/ (0 issues).
