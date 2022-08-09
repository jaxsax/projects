package main

import (
	"context"
	"log"
	"net/url"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
	"github.com/jessevdk/go-flags"
)

var (
	flagParser = flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)
	dbOptions  = &db.Options{}
	logOptions = &loggingOptions{}
)

type loggingOptions struct {
	DevelopmentLog bool `long:"pretty_logs" description:"use the nicer-to-look at development log" env:"PRETTY_LOGS"`
}

func main() {
	if _, err := flagParser.AddGroup("logging", "", logOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.AddGroup("db", "", dbOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}

	store, err := db.Setup(dbOptions)
	if err != nil {
		panic(err)
	}

	log.Printf("dsn=%s", dbOptions.URI)

	ctx := context.Background()
	allLinks, err := store.ListLinks(ctx)
	if err != nil {
		panic(err)
	}

	linksToUpdate := make([]*types.Link, 0)
	for _, link := range allLinks {
		linkURL, err := url.Parse(link.Link)
		if err != nil {
			panic(err)
		}

		path := linkURL.EscapedPath()

		if linkURL.RawQuery != "" {
			path += "?" + linkURL.RawQuery
		}

		if linkURL.Fragment != "" {
			path += "#" + linkURL.EscapedFragment()
		}

		link.Domain = linkURL.Hostname()
		link.Path = path

		linksToUpdate = append(linksToUpdate, link)
	}

	err = store.UpdateLinks(ctx, linksToUpdate)
	if err != nil {
		panic(err)
	}
}
