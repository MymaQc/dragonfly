# D36 — network chunk decode trusts third-party data

Branch: `chunk-network-decode-validate` (commit 388b9cc1, based on upstream/master 106fdc28)
Package: server/world/chunk

## Verdict
CONFIRMED (reproduced with a standalone decoder program).
chunk.NetworkDecode / NetworkDecodeBuffer accepted:
- a palette index size other than 0-6, 8 or 16 bits (only sizes > 32 were rejected);
- a block storage pointing to the previous one (0x7f), left as a nil storage -> nil deref on the first block read;
- indices pointing past the palette -> index out of range on the first block read;
- a palette entry count allocated before being read (up to gigabytes for a made-up count).

## Who can trigger it
Third-party / modified-server data: the server a client or proxy reads. Not a vanilla player on a Dragonfly server.

## Reproduction
dragonfly-bugs/repro/d36 (standalone): crafts network sub-chunk bytes for each case and decodes them
(single-sub-chunk range to keep the biome section tiny).
- before (wt/upstream-master, before.txt): nil back-reference -> PANIC; out-of-palette -> PANIC;
  invalid size 7 -> accepted; huge palette -> allocates then fails late.
- after  (wt/chunk-network-decode-validate, after.txt): each case rejected with a clear error; baseline still decodes.

## Files
- server/world/chunk/decode.go: reject invalid palette size; reject nil (back-reference) block storage; reject
  out-of-palette indices after decode.
- server/world/chunk/encoding.go: reject a network palette count larger than the bytes left before allocating.
- server/world/chunk/palette.go: paletteSize.valid().
- server/world/chunk/paletted_storage.go: PalettedStorage.withinPalette().

## Checks
go build, gofmt, go test ./server/world/chunk/ and ./server/world/ (round-trip tests pass), golangci-lint ./server/world/chunk/ (0 issues).
