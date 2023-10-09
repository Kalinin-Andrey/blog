package cli

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/order_event_full"
	"github.com/spf13/cobra"
	"github.com/wildberries-tech/wblogger"
)

// mpOrderCreationCollector ...
var defectCollector = &cobra.Command{
	Use:   "defect-collector",
	Short: "It is the defect-collector command.",
	Long:  `It is the defect-collector command: listen and count defects.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.defectCollector(cmd, args)
	},
}

func (app *App) defectCollector(cmd *cobra.Command, args []string) {
	metricName := "rating.defect-collector.defectCollector"
	var err error
	var event *order_event_full.OrderCreationEventFull
	// Process errors
	go func() {
		for err = range app.Integration.OrdoKafka.ChErr {
			wblogger.Error(app.ctx, metricName+" OrdoKafka.Consume error", err)
		}
	}()
	// Process results
	go func() {
		for event = range app.Integration.OrdoKafka.ChRes {
			if err := app.Domain.OrderEventFull.CreateNewEvent(app.ctx, event); err != nil {
				wblogger.Error(app.ctx, metricName+" OrderEventFull.CreateNewEvent error", err)
			}
		}
	}()

	if err := app.Integration.OrdoKafka.Consume(app.ctx); err != nil {
		wblogger.Error(app.ctx, metricName+" OrdoKafka.Consume exit with error", err)
	}
}

//func (app *App) defectCollector(cmd *cobra.Command, args []string) {
//	const metricName = "rating.mp-order-creation-collector.mpOrderCreationCollector"
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
//	err = mp_nats.PullSubscriptionProcessing(app.ctx, subs, app.defectHandler, config.Duration, 0)
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.MpNats.AddPullConsumerIfNotExists() error", err)
//		return
//	}
//}

//func (app *App) defectCollector(cmd *cobra.Command, args []string) {
//	const metricName = "rating.defect-collector.defectCollector"
//	config := app.config.DefectCollector
//	stat := app.Integration.OrdoNats.Conn().Stats()
//	status := app.Integration.OrdoNats.Conn().Status()
//	fmt.Println(stat)
//	fmt.Println(status)
//	sub, err := app.Integration.OrdoNats.Conn().QueueSubscribe(config.SubjectName, config.ConsumerGroupName, app.defectHandler)
//	if err != nil {
//		wblogger.Error(app.ctx, metricName+" Integration.OrdoNats.Conn().QueueSubscribe() error", err)
//		return
//	}
//	defer sub.Unsubscribe()
//}
//
//func (app *App) defectHandler(msg *nats.Msg) {
//	const metricName = "rating.defect-collector.defectHandler"
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
