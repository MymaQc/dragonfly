# D22 — saved firework entity whose Item is not a firework

Branch: `entity-firework-nbt-decode` (commit 821b2fed, based on upstream/master 106fdc28)
Package: server/entity

## Verdict
CONFIRMED on unmodified master. fireworkType.DecodeNBT asserts item.MapNBT(m,"Item").Item().(item.Firework);
a saved firework entity whose Item is missing or holds another item panics while the entity is read.

## Who can trigger it
World data only: a firework entity in a saved world whose Item field is not a firework (comparable to unknown
entity type handling, D8). Crashes the server as the chunk loads.

## Reproduction
dragonfly-bugs/repro/firework-entity (standalone). before.txt (wt/upstream-master) vs after.txt (wt/entity-firework-nbt-decode).

## Files
server/entity/firework.go

## Checks
go build, gofmt, go test ./server/entity/, golangci-lint ./server/entity/ (0 issues).
