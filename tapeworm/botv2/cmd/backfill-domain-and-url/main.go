package main

import (
	"context"
	"net/url"

	"github.com/go-logr/logr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type App struct {
	store     *db.Store
	dbOptions *db.Options
	logger    logr.Logger
}

func main() {
	app, err := initialize()
	if err != nil {
		panic(err)
	}

	app.logger.Info("info", "dsn", app.dbOptions.URI)

	ctx := context.Background()
	allLinks, err := app.store.ListLinks(ctx)
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

	err = app.store.UpdateLinks(ctx, linksToUpdate)
	if err != nil {
		panic(err)
	}
}
