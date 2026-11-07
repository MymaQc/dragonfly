# D19 — weather methods race the tick on the world's rand source

Branch: `world-weather-rand-race` (commit c0749533, based on upstream/master 106fdc28)
Package: server/world

## Verdict
CONFIRMED (go run -race: 49 DATA RACE warnings on master, 0 with the fix).
StopRaining and StopThundering (promoted onto World, callable from any goroutine) draw the next weather duration
from the World's own *math/rand/v2.Rand under the settings lock, while the ticking goroutine draws from the same
source without that lock (random block ticks, lightning). *rand.Rand is not safe for concurrent use; the race
detector flags concurrent access between weather.StopRaining and the tick goroutine's randUint4/tickLightning.

## Who can trigger it
Server / plugin code calling World.StartRaining/StopRaining/StopThundering from a goroutine other than the
world's, while the world ticks. Not a vanilla player.

## Reproduction
dragonfly-bugs/repro/d19 (standalone): loads a chunk, starts thundering, then spams Start/StopRaining from
another goroutine for 2s. Needs cgo/gcc for -race (available here).
  cd dragonfly-bugs/repro/d19 ; CGO_ENABLED=1 go run -race .
- before (wt/upstream-master): many "WARNING: DATA RACE" (before.txt; before_racecount.txt = 49)
- after  (wt/world-weather-rand-race): clean, prints "done" (after.txt; after_racecount.txt = 0)

## Files
server/world/weather.go (StopRaining, StopThundering draw the next duration from the goroutine-safe global
math/rand/v2 functions; same bounds; advanceWeather and the tick-goroutine uses of w.r are unchanged).

## Checks
go build, gofmt, go test ./server/world/, golangci-lint ./server/world/ (0 issues).
