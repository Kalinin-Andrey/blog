package cli

import (
	"context"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// calculator ...
var calculator = &cobra.Command{
	Use:   "calculator",
	Short: "It is the calculator command.",
	Long:  `It is the calculator command: listen and count new orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.calculator(cmd, args)
	},
}

func (app *App) calculator(cmd *cobra.Command, args []string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg := app.config.Calculator

	go func() {
		ticker := time.NewTicker(cfg.Duration)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				app.calculatorRun(app.ctx, cfg)
			}
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(app.ctx, shutdownTimeout)
	defer cancel()
	<-ctx.Done()

	return
}

func (app *App) calculatorRun(ctx context.Context, cfg *config.Calculator) {
	const metricName = "calculator.calculator"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.RatingProcessing.Process(ctx, cfg.RatingProcessingConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" RatingProcessing.Process() error", err)
			}
			return
		}(i)
	}

	return
}
