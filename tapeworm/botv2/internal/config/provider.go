package config

import (
	"sync"

	"github.com/google/wire"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/httpserver"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/telegrampoller"
	"github.com/jessevdk/go-flags"
)

var (
	flagParser = flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)
	once       sync.Once
	HTTP       = &httpserver.Options{}
	Telegram   = &telegrampoller.Options{}
	DB         = &db.Options{}
	Log        = &logging.Options{}

	m sync.Mutex
)

func initParser() {
	m.Lock()
	defer m.Unlock()

	once.Do(func() {
		if _, err := flagParser.AddGroup("http", "", HTTP); err != nil {
			panic(err)
		}

		if _, err := flagParser.AddGroup("telegram", "", Telegram); err != nil {
			panic(err)
		}

		if _, err := flagParser.AddGroup("db", "", DB); err != nil {
			panic(err)
		}

		if _, err := flagParser.AddGroup("logging", "", Log); err != nil {
			panic(err)
		}

		if _, err := flagParser.Parse(); err != nil {
			panic(err)
		}
	})
}

var ProviderSet = wire.NewSet(
	ProvideHTTP,
	ProvideDB,
	ProvideTelegram,
	ProvideLogging,
)

func ProvideHTTP() *httpserver.Options {
	initParser()

	return HTTP
}

func ProvideDB() *db.Options {
	initParser()

	return DB
}

func ProvideTelegram() *telegrampoller.Options {
	initParser()

	return Telegram
}

func ProvideLogging() *logging.Options {
	initParser()

	return Log
}
