package cli

import (
	"context"
	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
	"os"
	"os/signal"
	"syscall"
)

// filler-seller-id-of-nm-id ...
var fillerSellerIDOfNmID = &cobra.Command{
	Use:   "filler-seller-id-of-nm-id",
	Short: "It is the filler-seller-id-of-nm-id command.",
	Long:  `It is the filler-seller-id-of-nm-id command: fill seller_id of nm_id by new orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.fillerSellerIDOfNmID(cmd, args)
	},
}

func (app *App) fillerSellerIDOfNmID(cmd *cobra.Command, args []string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cfg := app.config.DefectCollector

	go func() {
		app.fillerSellerIDOfNmIDRun(app.ctx, cfg)
	}()

	<-stop

	ctx, cancel := context.WithTimeout(app.ctx, shutdownTimeout)
	defer cancel()
	<-ctx.Done()
	wblogger.Info(app.ctx, "filler-seller-id-of-nm-id: all done!")

	return
}

func (app *App) fillerSellerIDOfNmIDRun(ctx context.Context, cfg *config.DefectCollector) {
	const metricName = "rating.calculator.calculator"
	for i := byte(0); i < app.Infra.TsDB.Cluster.GetShardsNum(); i++ {
		go func(shardNum byte) {
			if err := app.Domain.OrderEventFull.FillSellerIDOfNmID(ctx, cfg.OrderEventFullConfig(), shardNum); err != nil {
				wblogger.Error(app.ctx, metricName+" RatingProcessing.Process() error", err)
			}
			return
		}(i)
	}

	return
}
