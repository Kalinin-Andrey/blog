package cli

import (
	"context"
	"errors"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// calculator ...
var defectProcessinger = &cobra.Command{
	Use:   "defect-processinger",
	Short: "It is the defect-processinger command.",
	Long:  `It is the defect-processinger command: listen and count new orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.defectProcessinger(cmd, args)
	},
}

func (app *App) defectProcessinger(cmd *cobra.Command, args []string) {
	stop := make(chan os.Signal, 1)
	defer close(stop)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg := app.config.DefectProcessinger

	go func() {
		app.defectProcessingerRun(app.ctx, stop, cfg)
	}()

	<-stop

	ctx, cancel := context.WithTimeout(app.ctx, shutdownTimeout)
	defer cancel()
	<-ctx.Done()

	return
}

func (app *App) defectProcessingerRun(ctx context.Context, stop chan os.Signal, cfg *config.DefectProcessinger) {
	const metricName = "defect-processinger.defectProcessingerRun"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {

			for {
				select {
				case <-stop:
					return
				default:
				}

				if err := app.Domain.DefectProcessing.Process(ctx, cfg.DefectProcessingConfig(), shardNum); err != nil {
					if errors.Is(err, apperror.ErrNotFound) {
						wblogger.Info(app.ctx, metricName+" all done!")
					} else {
						wblogger.Error(app.ctx, metricName+" DefectProcessing.Process() error", err)
					}
					time.Sleep(cfg.Duration)
				}
			}
			return
		}(i)
	}

	return
}
