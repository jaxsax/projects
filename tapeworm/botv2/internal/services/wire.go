package services

import (
	"github.com/google/wire"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/services/algolia"
	contentblock "github.com/jaxsax/projects/tapeworm/botv2/internal/services/content_block"
	dimcollector "github.com/jaxsax/projects/tapeworm/botv2/internal/services/dim_collector"
)

var Set = wire.NewSet(
	contentblock.New,
	algolia.NewForHN,
	dimcollector.New,
)
