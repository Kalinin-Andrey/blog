package cli

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/order_event"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
)

// mpOrderCreationCollector ...
var mpOrderCreationCollector = &cobra.Command{
	Use:   "mp-order-creation-collector",
	Short: "It is the mp-order-creation-collector command.",
	Long:  `It is the mp-order-creation-collector command: listen and count new orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.mpOrderCreationCollector(cmd, args)
	},
}

func (app *App) mpOrderCreationCollector(cmd *cobra.Command, args []string) {
	const metricName = "mp-order-creation-collector.mpOrderCreationCollector"

	var err error
	var event *order_event.OrderCreationEvent
	// Process errors
	go func() {
		for err = range app.Integration.MpKafka.NewOrders.ChErr {
			wblogger.Error(app.ctx, metricName+" OrdoKafka.Consume error", err)
		}
	}()
	// Process results
	go func() {
		for event = range app.Integration.MpKafka.NewOrders.ChRes {
			if err := app.Domain.OrderEvent.CreateOrderCreationEvent(app.ctx, event); err != nil {
				wblogger.Error(app.ctx, metricName+" OrderEvent.CreateOrderCreationEvent error", err)
			}
		}
	}()

	if err := app.Integration.MpKafka.NewOrders.Run(app.ctx); err != nil {
		wblogger.Error(app.ctx, metricName+" MpKafka.NewOrders.Run exit with error", err)
	}
}

// Для работы с NATS
//func (app *App) mpOrderCreationCollector(cmd *cobra.Command, args []string) {
//	const metricName = "mp-order-creation-collector.mpOrderCreationCollector"
//
//	config := app.config.MpOrderCreationCollector
//	isStreamExists, _, err := app.Integration.MpNats.IsStreamExists(config.StreamName)
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.IsStreamExists() error", err)
//		return
//	}
//	if isStreamExists == false {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.IsStreamExists() = false", err)
//		return
//	}
//
//	_, subs, delFunc, err := app.Integration.MpNats.AddPullConsumerIfNotExists(config.StreamName, &nats.ConsumerConfig{
//		Name:          config.ConsumerName,
//		Durable:       config.ConsumerName,
//		FilterSubject: config.SubjectName,
//	})
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.AddPullConsumerIfNotExists() error", err)
//		return
//	}
//	defer func() {
//		if err = delFunc(); err != nil {
//			wblogger.Error(app.ctx, metricName+" delFunc() error", err)
//		}
//	}()
//
//	err = mp_nats.PullSubscriptionProcessing(app.ctx, subs, app.mpNewOrderHandler, config.Duration, 0)
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.AddPullConsumerIfNotExists() error", err)
//		return
//	}
//}
//
//func (app *App) mpNewOrderHandler(msg *nats.Msg) {
//	const metricName = "mp-order-creation-collector.mpNewOrderHandler"
//	newOrder := &mp_message.NewOrder{}
//	if err := newOrder.UnmarshalBinary(msg.Data); err != nil {
//		wblogger.Error(app.ctx, metricName+" newOrder.UnmarshalBinary() error:", err)
//		return
//	}
//
//	err := app.Domain.OrderEvent.CreateOrderCreationEvent(app.ctx, mp_message.NewOrderProto2OrderCreationEvent(newOrder))
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" OrderEvent.CreateOrderCreationEvent() error", err)
//	}
//	return
//}
