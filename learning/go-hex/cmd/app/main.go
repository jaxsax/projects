package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jaxsax/projects/learning/go-hex/internal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	registry, err := internal.NewRegistry(&internal.Config{
		SQLPath: "./store.db?_journal_mode=WAL",
	})
	if err != nil {
		log.Fatalf("registry create err: %v", err)
	}

	app := internal.NewApplication(registry)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("err: %v", err)
		}
	}()

	<-ch
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Printf("Shutting down..")
	if err := app.Stop(shutdownCtx); err != nil {
		log.Fatalf("shutdown err: %v", err)
	}
}
