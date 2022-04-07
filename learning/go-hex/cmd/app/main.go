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

	app := internal.NewApplication()

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
