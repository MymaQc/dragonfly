# D28 — smithing table accepts empty slots

Branch: `session-smithing-empty-slot` (commit a98e202f, based on upstream/master 106fdc28)
Package: server/session

## Verdict
CONFIRMED at the handler level (root cause reproduced).
handleSmithing checks the input, material and template slots with matchingStacks, which returns true whenever
either stack is empty (item.Stack.Comparable short-circuits on an empty stack). A smithing recipe always has
three non-empty ingredients, so an empty slot is never a real match there. A modified client can send a
SmithingTransform request for a valid recipe with the table slots empty: every check passed and the product
was created without consuming anything.

## Who can trigger it
Modified client only: it must send an ItemStackRequest (CraftRecipe) with a valid smithing recipe network id
while leaving the table slots empty. A vanilla client always has the items in the slots.

## Reproduction
Root cause (matchingStacks) shown by a session-package test, kept OUTSIDE the fix commit:
  evidence/session-smithing-empty-slot/matchingstacks_test.go.txt
Run from the master worktree: copy it to server/session/, `go test ./server/session/ -run TestD28 -v` (see before.txt),
then delete it. It shows matchingStacks(empty, NetheriteIngot) == true, so handleSmithing would accept the empty slot.
A full client-level run needs a live modified client; the fix (after.txt) guards each slot with Stack.Empty().

## Files
server/session/handler_smithing.go (reject empty input/material/template before matchingStacks).

## Checks
go build, gofmt, go test ./server/session/, golangci-lint ./server/session/ (0 issues).
