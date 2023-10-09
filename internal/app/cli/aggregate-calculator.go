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

// aggregateCalculator ...
var aggregateCalculator = &cobra.Command{
	Use:   "aggregate-calculator",
	Short: "It is the aggregate-calculator command.",
	Long:  `It is the aggregate-calculator command: calculates and saves the aggregated metrics for the rating.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.aggregateCalculator(cmd, args)
	},
}

func (app *App) aggregateCalculator(cmd *cobra.Command, args []string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg := app.config.AggregateCalculator

	go func() {
		app.aggregateCalculatorSinc(app.ctx, cfg)
		ticker := time.NewTicker(cfg.SincDuration)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				app.aggregateCalculatorSinc(app.ctx, cfg)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(cfg.Duration)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				app.aggregateCalculatorRun(app.ctx, cfg)
			}
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(app.ctx, shutdownTimeout)
	defer cancel()
	<-ctx.Done()

	return
}

func (app *App) aggregateCalculatorRun(ctx context.Context, cfg *config.AggregateCalculator) {
	const metricName = "rating.aggregate-calculator.aggregateCalculatorRun"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.AggregateForRatingProcessing.Process(ctx, cfg.AggregateForRatingProcessingConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" AggregateForRatingProcessing.Process() error", err)
			}
			return
		}(i)
	}

	return
}

func (app *App) aggregateCalculatorSinc(ctx context.Context, cfg *config.AggregateCalculator) {
	const metricName = "rating.aggregate-calculator.aggregateCalculatorSinc"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.AggregateForRatingProcessing.SyncForProcessingList(ctx, cfg.AggregateForRatingProcessingConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" AggregateForRatingProcessing.SyncForProcessingList() error", err)
			}
			return
		}(i)
	}

	return
}
