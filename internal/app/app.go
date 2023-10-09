package app

import (
	"context"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
	"github.com/wildberries-tech/wblogger"
	"log"

	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/redis"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb_cluster"
)

type App struct {
	config *config.AppConfig
	Infra  *infrastructure.Infrastructure
	Domain *Domain
}

type Domain struct {
	Blog *blog.Service
}

// New func is a constructor for the App
func New(ctx context.Context, cfg *config.Configuration) *App {
	log.Println("Core app is starting...")
	wblogger.Debug(ctx, "app.New()")
	log.Println("infrastructure start create...")
	infr, err := infrastructure.New(ctx, cfg.App.InfraAppConfig(), cfg.Infra)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")

	app := &App{
		config: cfg.App,
		Infra:  infr,
	}

	app.SetupServices()

	return app
}

func (app *App) SetupServices() {
	app.Domain = &Domain{
		Blog: blog.NewService(redis.NewBlogReplicaSet(app.Infra.Redis), tsdb_cluster.NewBlogReplicaSet(app.Infra.TsDB)),
	}
}

func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	return app.Infra.Close()
}
