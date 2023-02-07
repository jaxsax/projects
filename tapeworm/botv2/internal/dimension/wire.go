package dimension

import "github.com/google/wire"

var Set = wire.NewSet(
	NewHNCollector,
)
