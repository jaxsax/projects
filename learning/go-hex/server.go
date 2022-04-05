package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Application struct {
	httpServer *http.Server
}

func NewApplication() *Application {
	app := &Application{}

	g := gin.Default()
	app.setupRouter(g)

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: g,
	}

	app.httpServer = httpServer

	return app
}

func (a *Application) Run() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}

func (a *Application) setupRouter(g *gin.Engine) {
	g.GET("/upper", func(ctx *gin.Context) {
		q := ctx.Query("q")
		q = strings.ToLower(q)
		_, _ = ctx.Writer.Write([]byte(q))
	})
}
