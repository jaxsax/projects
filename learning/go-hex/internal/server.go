package internal

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Application struct {
	httpServer *http.Server

	Router *gin.Engine
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
	app.Router = g

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
	log.Printf("Stopping http")
	return a.httpServer.Shutdown(ctx)
}

func (a *Application) setupRouter(g *gin.Engine) {
	g.GET("/upper", func(ctx *gin.Context) {
		q := ctx.Query("q")
		q = strings.ToUpper(q)
		_, _ = ctx.Writer.Write([]byte(q))
	})
}
