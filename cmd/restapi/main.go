package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/wildberries-tech/wblogger"
	_ "go.uber.org/automaxprocs"

	"github.com/Kalinin-Andrey/blog/internal/app"
	"github.com/Kalinin-Andrey/blog/internal/app/restapi"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
)

var Version = "0.0.0"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		wblogger.Flush()
	}()

	conf, err := config.Get()
	if err != nil {
		log.Fatal(errors.Wrap(err, "load conf error"))
	}
	coreApp := app.New(ctx, conf)
	restAPI := restapi.New(coreApp, conf.App, conf.API)

	done := make(chan os.Signal, 1)
	go func() {
		if err = restAPI.Run(ctx); err != nil {
			log.Fatal(errors.Wrap(err, "start of rest api error"))
		}
	}()

	defer func() {
		if err = restAPI.Stop(); err != nil {
			wblogger.Error(ctx, "API Shutdown error", err)
		}
	}()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
