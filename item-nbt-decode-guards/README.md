# D20 — item NBT decode panics

Branch: `item-nbt-decode-guards` (commit 6e17fa86, based on upstream/master 106fdc28)
Package: server/item (grouped: written_book, firework, firework_explosion; firework_star covered via firework_explosion)

## Verdict
CONFIRMED on unmodified master. Reproduced with repro_main.go (rows D20).

- Written book page that is not a compound, or whose "text" is not a string -> panic.
- Firework rocket whose Explosions holds a non-compound entry -> panic.
- Firework explosion whose FireworkType/FireworkFlicker/FireworkTrail is missing or wrong type, or whose
  FireworkType indexes past FireworkShapes() -> panic. FireworkStar reuses FireworkExplosion.DecodeNBT, so it is fixed too.

## Who can trigger it
Client-triggerable: item NBT arrives through stackToItem (server/session/player.go) whenever a client sends an
item stack (creative inventory, item stack request, mob equipment). A crafted stack crashes the server.

## Reproduction
dragonfly-bugs/repro/nbt, rows D20. Switch the `replace` to wt/upstream-master (before.txt) or
wt/item-nbt-decode-guards (after.txt) and `go run .`.

## Files
server/item/written_book.go, firework.go, firework_explosion.go

## Checks
go build, gofmt, go test ./server/item/, golangci-lint ./server/item/ (0 issues).

Note: server/item cannot import internal/nbtconv (import cycle), so the firework explosion fields use inline comma-ok reads rather than nbtconv helpers.
