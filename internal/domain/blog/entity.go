package blog

import "time"

const ()

type Blog struct {
	Timestamp           time.Time `json:"timestamp"`           // время рассчёта
	SellerName          string    `json:"sellerName"`          // 2	Наименование поставщика
	Rating              float64   `json:"rating"`              // 9	Рейтинг
	RatioDelivered      float64   `json:"ratioDelivered"`      // 6	Доля доставленных
	RatioDefected       float64   `json:"ratioDefected"`       // 15	Доля брака
	BuyerRatingWeight   float64   `json:"buyerRatingWeight"`   // 12	Влияние оценок на общий рейтинг продавца
	AvgBuyerRating      float64   `json:"avgBuyerRating"`      // 10	Средняя оценка за товары
	SellerOldId         uint      `json:"sellerOldId"`         // 1	ИД (SupplierOldId [int])
	NbDelivered         uint      `json:"nbDelivered"`         // 4	Кол-во доставленных на склад
	NbInDelivery        uint      `json:"nbInDelivery"`        // 5	Кол-во в ожидании (в доставке)
	NbOrdersMarketplace uint      `json:"nbOrdersMarketplace"` // 3	Кол-во заказанных ШК srid
	NbBuyerRatings      uint      `json:"nbBuyerRatings"`      // 11	Количество оценок за товары
	NbDefected          uint      `json:"nbDefected"`          // 13	Кол-во брака
	NbOrdersTotal       uint      `json:"nbOrdersTotal"`       // 14	Общее количество заказов (из ОРДО)
}

func (e *Blog) Validate() error {
	return nil
}
