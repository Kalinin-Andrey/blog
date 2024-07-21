package cli

import (
	"blog/internal/pkg/config"
	"context"
	"errors"
	prometheus_utils "github.com/minipkg/prometheus-utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net/http"

	"blog/internal/app"
)

const (
	metricsSuccess = "true"
	metricsFail    = "false"
)

// App is the application for CLI app
type App struct {
	*app.App
	ctx           context.Context
	config        *config.CliConfig
	apiConfig     *config.API
	logger        *zap.Logger
	rootCmd       *cobra.Command
	serverMetrics *fasthttp.Server
	serverProbes  *fasthttp.Server
}

var CliApp *App

// New func is a constructor for the App
func New(ctx context.Context, appName string, coreApp *app.App, apicfg *config.API, cfg *config.CliConfig) *App {
	CliApp = &App{
		App:       coreApp,
		ctx:       ctx,
		config:    cfg,
		apiConfig: apicfg,
		logger:    coreApp.Infra.Logger,
		rootCmd: &cobra.Command{
			Use:   "cli",
			Short: "This is the short description.",
			Long:  `This is the long description.`,
		},
		serverMetrics: &fasthttp.Server{
			Name:            appName,
			ReadTimeout:     apicfg.Metrics.ReadTimeout,
			WriteTimeout:    apicfg.Metrics.WriteTimeout,
			IdleTimeout:     apicfg.Metrics.IdleTimeout,
			CloseOnShutdown: true,
		},
		serverProbes: &fasthttp.Server{
			Name:            appName,
			ReadTimeout:     apicfg.Probes.ReadTimeout,
			WriteTimeout:    apicfg.Probes.WriteTimeout,
			IdleTimeout:     apicfg.Probes.IdleTimeout,
			CloseOnShutdown: true,
		},
	}
	CliApp.init()
	return CliApp
}

func (app *App) init() {
	app.rootCmd.AddCommand()
	app.buildHandler()
}

func (a *App) buildHandler() {
	rp := routing.New()
	rp.Get("/live", LiveHandler)
	rp.Get("/ready", LiveHandler)
	a.serverProbes.Handler = rp.HandleRequest

	rm := routing.New()
	rm.Get("/metrics", prometheus_utils.GetFasthttpRoutingHandler())
	a.serverMetrics.Handler = rm.HandleRequest
}

// Run is func to run the App
func (app *App) Run() error {
	if err := app.App.Run(); err != nil {
		return err
	}
	go func() {
		app.logger.Info("metrics listen on " + app.apiConfig.Metrics.Addr)
		if err := app.serverMetrics.ListenAndServe(app.apiConfig.Metrics.Addr); err != nil {
			app.logger.Error("serverMetrics.ListenAndServe error", zap.Error(err))
		}
	}()
	go func() {
		app.logger.Info("probes listen on " + app.apiConfig.Probes.Addr)
		if err := app.serverProbes.ListenAndServe(app.apiConfig.Probes.Addr); err != nil {
			app.logger.Error("serverProbes.ListenAndServe error", zap.Error(err))
		}
	}()
	app.logger.Info("cli app is starting...")
	return app.rootCmd.Execute()
}

func (app *App) Stop() error {
	app.logger.Info("Cli-Shutdown")
	return errors.Join(
		app.serverMetrics.Shutdown(),
		app.serverProbes.Shutdown(),
		app.App.Stop(),
	)
}

func LiveHandler(rctx *routing.Context) error {
	rctx.SetStatusCode(http.StatusNoContent)
	return nil
}
