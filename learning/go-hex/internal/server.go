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
	registry   *Registry

	Router *gin.Engine
}

func NewApplication(reg *Registry) *Application {
	app := &Application{}

	g := gin.Default()
	app.setupRouter(g)

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: g,
	}

	app.httpServer = httpServer
	app.Router = g
	app.registry = reg

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

	g.GET("/pets", func(ctx *gin.Context) {
		type Pet struct {
			ID   uint64 `json:"id"`
			Name string `json:"name"`
		}

		type PetsResponse struct {
			Pets []Pet `json:"pets"`
		}

		rows, err := a.registry.DB.Queryx("SELECT * FROM pets")
		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError, gin.H{

					"type":   "server-error",
					"title":  "Server error",
					"status": http.StatusInternalServerError,
					"detail": err.Error(),
				},
			)
			return
		}

		pets := make([]Pet, 0)
		for rows.Next() {
			var pet Pet
			if err := rows.StructScan(&pet); err != nil {
				ctx.JSON(
					http.StatusInternalServerError, gin.H{

						"type":   "server-error",
						"title":  "Server error",
						"status": http.StatusInternalServerError,
						"detail": err.Error(),
					},
				)
				return
			}

			pets = append(pets, pet)
		}

		ctx.JSON(http.StatusOK, &PetsResponse{
			Pets: pets,
		})
	})
}
