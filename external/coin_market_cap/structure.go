package coin_market_cap

import "time"

type PriceConversionResponse struct {
	Status interface{}                 `json:"status"`
	Data   PriceConversionDataResponse `json:"data"`
}

type PriceConversionDataResponse struct {
	Id          int       `json:"id"`
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Amount      int       `json:"amount"`
	LastUpdated time.Time `json:"last_updated"`
	Quote       Quote     `json:"quote"`
}
type Quote struct {
	USD USD `json:"usd"`
}

type USD struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}
