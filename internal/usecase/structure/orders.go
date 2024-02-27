package structure

type FilterOrders struct {
	Email *string `json:"email"`
}

type ApiOrderResp struct {
	Status  int                  `json:"status"`
	Data    ApiOrderDataResp     `json:"data"`
	Message *ApiOrderMessageResp `json:"message"`
}

type ApiOrderMessageResp struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

type ApiOrderDataResp struct {
	Orders []*ApiOrderItemResp `json:"orders"`
	Total  int                 `json:"total"`
}

type ApiOrderItemResp struct {
	Id                 string                   `json:"id"`
	Amount             string                   `json:"amount"`
	Status             int                      `json:"status"`
	PayType            string                   `json:"pay_type"`
	ShippingFirstname  string                   `json:"shipping_firstname"`
	ShippingAddress1   string                   `json:"shipping_address1"`
	ShippingAddress2   string                   `json:"shipping_address2"`
	ShippingCity       string                   `json:"shipping_city"`
	ShippingRegion     string                   `json:"shipping_region"`
	ShippingPostalCode string                   `json:"shipping_postal_code"`
	ShippingCountry    string                   `json:"shipping_country"`
	Email              string                   `json:"email"`
	EvmMasterAddress   string                   `json:"evm_master_address"`
	ShippingDate       string                   `json:"shipping_date"`
	CouponCode         string                   `json:"coupon_code"`
	DiscountAmount     string                   `json:"discount_amount"`
	CouponType         int                      `json:"coupon_type"`
	CouponValue        string                   `json:"coupon_value"`
	CreatedAt          string                   `json:"created_at"`
	Details            []ApiOrderItemDetailResp `json:"details"`
}

type ApiOrderItemDetailResp struct {
	ProductId       string `json:"product_id"`
	ProductPrice    string `json:"product_price"`
	ProductName     string `json:"product_name"`
	ProductOptionId string `json:"product_option_id"`
	CustomLogo      string `json:"custom_logo"`
	Quantity        int    `json:"quantity"`
	Amount          string `json:"amount"`
	EthRate         string `json:"eth_rate"`
	EthRateTime     string `json:"eth_rate_time"`
}

type OrderStatusChan struct {
	OrderID string
	PayType string
	Err     error
	Status  int
}
