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

// feedbackCollector ...
var feedbackCollector = &cobra.Command{
	Use:   "feedback-collector",
	Short: "It is the feedback-collector command.",
	Long:  `It is the feedback-collector command: get and save new grades from feedbacks.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.feedbackCollector(cmd, args)
	},
}

func (app *App) feedbackCollector(cmd *cobra.Command, args []string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg := app.config.FeedbackCollector

	go func() {
		app.feedbackCollectorSinc(app.ctx, cfg)
		ticker := time.NewTicker(cfg.SincDuration)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				app.feedbackCollectorSinc(app.ctx, cfg)
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
				app.feedbackCollectorRun(app.ctx, cfg)
			}
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(app.ctx, shutdownTimeout)
	defer cancel()
	<-ctx.Done()

	return
}

func (app *App) feedbackCollectorRun(ctx context.Context, cfg *config.FeedbackCollector) {
	const metricName = "rating.feedback-collector.feedbackCollectorRun"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.FeedbackProcessing.Process(ctx, cfg.FeedbackProcessingConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" FeedbackProcessing.Process() error", err)
			}
			return
		}(i)
	}

	return
}

func (app *App) feedbackCollectorSinc(ctx context.Context, cfg *config.FeedbackCollector) {
	const metricName = "rating.feedback-collector.feedbackCollectorSinc"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.FeedbackProcessing.SyncForProcessingList(ctx, cfg.FeedbackProcessingConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" FeedbackProcessing.SyncForProcessingList() error", err)
			}
			return
		}(i)
	}

	return
}
