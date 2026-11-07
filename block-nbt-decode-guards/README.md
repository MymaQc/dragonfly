# D26 / D30 / D31 — block entity NBT decode panics

Branch: `block-nbt-decode-guards` (commit 52c0332c, based on upstream/master 106fdc28)
Package: server/block (grouped: beacon, bed, banner, decorated_pot, skull)

## Verdict
CONFIRMED on unmodified master. Reproduced with repro_main.go.

Covers three divergence rows, all in server/block:
- D26 beacon: Primary/Secondary that is a non-lasting effect id -> panic on effect.LastingType assertion.
- D30 bed: color > 15 -> index out of range on item.Colours().
- D31 banner: a pattern layer that is not a compound, or an unknown pattern id -> panic;
  decorated pot: a sherd that is not a string, an unknown item, or an item that is not a PotDecoration -> panic;
  skull: SkullType 7..254 -> skull.String() panics when the block state is later encoded.

## Who can trigger it
- World data only for all except the items that can carry this NBT.
- Client-triggerable for the item forms: banner items and written/skull/pot item NBT arrive through stackToItem
  (server/session/player.go) when a client sends an item stack (creative inventory / item stack request / mob equipment).
  Bed and beacon are block-entity (world data) only.

## Reproduction
repro module: dragonfly-bugs/repro/nbt (standalone, calls each DecodeNBT with crafted maps).
  cd dragonfly-bugs/repro/nbt
  # before: replace -> wt/upstream-master ; go run .   (see before.txt, rows D26/D30/D31)
  # after:  replace -> wt/block-nbt-decode-guards ; go run .  (see after.txt)

## Files
- server/block/beacon.go, bed.go, banner.go, decorated_pot.go, skull.go

## Checks
go build, gofmt, go test ./server/block/, golangci-lint ./server/block/ (0 issues), go generate ./server/block/ (no diff).
