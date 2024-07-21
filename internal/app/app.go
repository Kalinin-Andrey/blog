package app

import (
	"blog/internal/integration"
	"blog/internal/pkg/config"
	"context"
	"errors"
	"log"

	"blog/internal/domain/blog"
	"blog/internal/infrastructure"
	"blog/internal/infrastructure/repository/tsdb_cluster"
)

type App struct {
	config      *config.AppConfig
	Infra       *infrastructure.Infrastructure
	Integration *integration.Integration
	Domain      *Domain
}

type Domain struct {
	Blog *blog.Service
}

// New func is a constructor for the App
func New(ctx context.Context, cfg *config.Configuration) *App {
	log.Println("Core app is starting...")
	log.Println("infrastructure start create...")
	infr, err := infrastructure.New(ctx, cfg.App.InfraAppConfig(), cfg.Infra)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")

	log.Println("integration start create...")
	integr, err := integration.New(&integration.AppConfig{
		NameSpace:   cfg.App.NameSpace,
		Subsystem:   cfg.App.Name,
		Service:     cfg.App.Service,
		Environment: cfg.App.Environment,
	}, cfg.Integration)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")

	app := &App{
		config:      cfg.App,
		Infra:       infr,
		Integration: integr,
	}

	app.SetupServices()

	return app
}

func (app *App) SetupServices() {
	app.Domain = &Domain{
		Blog: blog.NewService(tsdb_cluster.NewBlogReplicaSet(app.Infra.TsDB)),
	}
}

func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	return errors.Join(
		app.Integration.Close(),
		app.Infra.Close(),
	)
}
