package cli

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/order_event"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
)

// mpOrderChangeCollector ...
var mpOrderChangeCollector = &cobra.Command{
	Use:   "mp-order-change-collector",
	Short: "It is the mp-order-change-collector command.",
	Long:  `It is the mp-order-change-collector command: listen and count changing of state orders.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.mpOrderChangeCollector(cmd, args)
	},
}

func (app *App) mpOrderChangeCollector(cmd *cobra.Command, args []string) {
	const metricName = "mp-order-change-collector.mpOrderChangeCollector"

	//go app.mpOrderChangeCollector_ScannedOrders()
	go app.mpOrderChangeCollector_MPStates()
	//go app.mpOrderChangeCollector_OrdoStates()
}

func (app *App) mpOrderChangeCollector_ScannedOrders() {
	const metricName = "mp-order-change-collector.mpOrderChangeCollector_ScannedOrders"

	var err error
	var event *order_event.OrderChangeEvent
	// Process errors
	go func() {
		for err = range app.Integration.MpKafka.ScannedOrders.ChErr {
			wblogger.Error(app.ctx, metricName+" MpKafka.ScannedOrders.Consume error", err)
		}
	}()
	// Process results
	go func() {
		for event = range app.Integration.MpKafka.ScannedOrders.ChRes {
			if err := app.Domain.OrderEvent.CreateOrderChangeEvent(app.ctx, event); err != nil {
				wblogger.Error(app.ctx, metricName+" OrderEvent.CreateOrderChangeEvent error", err)
			}
		}
	}()

	if err := app.Integration.MpKafka.ScannedOrders.Run(app.ctx); err != nil {
		wblogger.Error(app.ctx, metricName+" MpKafka.ScannedOrders.Run exit with error", err)
	}
}

func (app *App) mpOrderChangeCollector_MPStates() {
	const metricName = "mp-order-change-collector.mpOrderChangeCollector_MPStates"

	var err error
	var event *order_event.OrderChangeEvent
	// Process errors
	go func() {
		for err = range app.Integration.MpKafka.MPStates.ChErr {
			wblogger.Error(app.ctx, metricName+" MpKafka.MPStates.Consume error", err)
		}
	}()
	// Process results
	go func() {
		for event = range app.Integration.MpKafka.MPStates.ChRes {
			if err := app.Domain.OrderEvent.CreateOrderChangeEvent(app.ctx, event); err != nil {
				wblogger.Error(app.ctx, metricName+" OrderEvent.CreateOrderChangeEvent error", err)
			}
		}
	}()

	if err := app.Integration.MpKafka.MPStates.Run(app.ctx); err != nil {
		wblogger.Error(app.ctx, metricName+" MpKafka.MPStates.Run exit with error", err)
	}
}

func (app *App) mpOrderChangeCollector_OrdoStates() {
	const metricName = "mp-order-change-collector.mpOrderChangeCollector_MPStates"

	var err error
	var event *order_event.OrderChangeEvent
	// Process errors
	go func() {
		for err = range app.Integration.MpKafka.OrdoStates.ChErr {
			wblogger.Error(app.ctx, metricName+" MpKafka.OrdoStates.Consume error", err)
		}
	}()
	// Process results
	go func() {
		for event = range app.Integration.MpKafka.OrdoStates.ChRes {
			if err := app.Domain.OrderEvent.CreateOrderChangeEvent(app.ctx, event); err != nil {
				wblogger.Error(app.ctx, metricName+" OrderEvent.CreateOrderChangeEvent error", err)
			}
		}
	}()

	if err := app.Integration.MpKafka.OrdoStates.Run(app.ctx); err != nil {
		wblogger.Error(app.ctx, metricName+" MpKafka.OrdoStates.Run exit with error", err)
	}
}

//// Для работы с NATS
//func (app *App) mpOrderChangeCollector(cmd *cobra.Command, args []string) {
//	const metricName = "mp-order-change-collector.mpOrderChangeCollector"
//
//	config := app.config.MpOrderChangeCollector
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
//	err = mp_nats.PullSubscriptionProcessing(app.ctx, subs, app.mpStateOrderHandler, config.Duration, 0)
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.AddPullConsumerIfNotExists() error", err)
//		return
//	}
//}
//
//func (app *App) mpStateOrderHandler(msg *nats.Msg) {
//	const metricName = "rating.mp-order-change-collector.mpStateOrderHandler"
//	orderStates := &mp_message.PurchaseOrderStates{}
//	if err := orderStates.UnmarshalBinary(msg.Data); err != nil {
//		wblogger.Error(app.ctx, metricName+" orderStates.UnmarshalBinary() error:", err)
//		return
//	}
//
//	for _, orderState := range orderStates.List {
//		err := app.Domain.OrderEvent.CreateOrderChangeEvent(app.ctx, mp_message.PurchaseOrderStateProto2OrderChangeEvent(orderState))
//		if err != nil {
//			wblogger.Error(app.ctx, metricName+" OrderEvent.CreateChanged() error", err)
//		}
//	}
//	return
//}
