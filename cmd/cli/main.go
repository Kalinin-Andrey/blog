package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "go.uber.org/automaxprocs"

	"blog/internal/app"
	"blog/internal/app/cli"
	"blog/internal/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	conf, err := config.Get()
	if err != nil {
		log.Fatalf("Load conf error: %w", err)
	}
	cliApp := cli.New(ctx, conf.App.Name, app.New(ctx, conf), conf.API, conf.Cli)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = cliApp.Run(); err != nil {
			log.Fatalf("cliApp.Run() error: %w", err)
			os.Exit(1)
		}
	}()

	defer func() {
		if err = cliApp.Stop(); err != nil {
			log.Printf("cliApp.Stop() error: %w", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
	wg.Wait()
}
