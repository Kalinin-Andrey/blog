package cli

import (
	"context"
	"errors"
	"net/http"
	"time"

	prometheus_utils "github.com/minipkg/prometheus-utils"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/wildberries-tech/wblogger"

	"github.com/Kalinin-Andrey/blog/internal/pkg/config"

	"github.com/Kalinin-Andrey/blog/internal/app"
)

const (
	shutdownTimeout = 3 * time.Second
)

// App is the application for CLI app
type App struct {
	*app.App
	ctx           context.Context
	config        *config.CliConfig
	apiConfig     *config.API
	rootCmd       *cobra.Command
	errors        []error
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
	app.rootCmd.AddCommand(mpOrderCreationCollector,
		mpOrderChangeCollector,
		calculator,
		feedbackCollector,
		defectCollector,
		fillerSellerIDOfNmID,
		defectProcessinger,
		aggregateCalculator,
	)
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
		wblogger.Info(context.Background(), "metrics listen on "+app.apiConfig.Metrics.Addr)
		if err := app.serverMetrics.ListenAndServe(app.apiConfig.Metrics.Addr); err != nil {
			wblogger.Error(app.ctx, "serverMetrics.ListenAndServe error", err)
			wblogger.Flush()
		}
	}()
	go func() {
		wblogger.Info(context.Background(), "probes listen on "+app.apiConfig.Probes.Addr)
		if err := app.serverProbes.ListenAndServe(app.apiConfig.Probes.Addr); err != nil {
			wblogger.Error(app.ctx, "serverProbes.ListenAndServe error", err)
			wblogger.Flush()
		}
	}()
	wblogger.Info(context.Background(), "cli app is starting...")
	return app.rootCmd.Execute()
}

func (app *App) Stop() error {
	var err error
	var errs []error

	wblogger.Info(context.Background(), "Cli-Shutdown")
	time.Sleep(time.Second * 10)

	if err = app.serverMetrics.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err = app.serverProbes.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err = app.App.Stop(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func LiveHandler(rctx *routing.Context) error {
	rctx.SetStatusCode(http.StatusNoContent)
	return nil
}
