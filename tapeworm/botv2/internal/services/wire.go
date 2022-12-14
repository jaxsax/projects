package services

import (
	"github.com/google/wire"
	contentblock "github.com/jaxsax/projects/tapeworm/botv2/internal/services/content_block"
)

var Set = wire.NewSet(
	contentblock.New,
)
