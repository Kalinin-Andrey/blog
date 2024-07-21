package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "go.uber.org/automaxprocs"

	"blog/internal/app"
	"blog/internal/app/restapi"
	"blog/internal/pkg/config"
)

var Version = "0.0.0"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	conf, err := config.Get()
	if err != nil {
		log.Fatalf("Load conf error: %w", err)
	}
	restAPI := restapi.New(app.New(ctx, conf), conf.App, conf.API)

	go func() {
		if err = restAPI.Run(ctx); err != nil {
			log.Fatalf("restAPI.Run(ctx) error: %w", err)
			os.Exit(1)
		}
	}()

	defer func() {
		if err = restAPI.Stop(); err != nil {
			log.Printf("restAPI.Stop() error: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
}
