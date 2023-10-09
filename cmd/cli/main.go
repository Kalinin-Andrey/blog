package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wildberries-tech/wblogger"
	_ "go.uber.org/automaxprocs"

	"github.com/Kalinin-Andrey/blog/internal/app"
	"github.com/Kalinin-Andrey/blog/internal/app/cli"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
)

func main() {
	ctx := context.Background()
	conf, err := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	cliApp := cli.New(ctx, conf.App.Name, app.New(ctx, conf), conf.API, conf.Cli)

	done := make(chan os.Signal, 1)
	go func() {
		if err = cliApp.Run(); err != nil {
			log.Fatalf("Error while cli application is running: %s", err.Error())
			os.Exit(1)
		}
	}()

	defer func() {
		if err = cliApp.Stop(); err != nil {
			wblogger.Error(ctx, "API Shutdown error", err)
		}
	}()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

}
