package internal

import (
	"github.com/google/wire"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/config"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

var CommonSet = wire.NewSet(
	config.ProviderSet,
	logging.ProviderSet,
	db.Setup,
)
